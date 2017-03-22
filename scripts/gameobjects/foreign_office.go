gameobject "foreign_office" inherit "office"
<<
	local 
	<<
		function eraseDiplomacyPoints()
			send("rendUIManager", "SetOfficeInt", SELF, "Diplomacy_Points_Empire"    , 0)
			send("rendUIManager", "SetOfficeInt", SELF, "Diplomacy_Points_Stahlmark" , 0)
			send("rendUIManager", "SetOfficeInt", SELF, "Diplomacy_Points_Novorus"   , 0)
			send("rendUIManager", "SetOfficeInt", SELF, "Diplomacy_Points_Republique", 0)
			send("rendUIManager", "SetOfficeInt", SELF, "Diplomacy_Points_Bandits"   , 0)
			state.diplomacyPoints.Stahlmark = 0
			state.diplomacyPoints.Novorus = 0
			state.diplomacyPoints.Republique = 0
			state.diplomacyPoints.Empire = 0
			state.diplomacyPoints.Bandits = 0
		end

		function changeStanding( faction, mission)
			--printl("office", mission.name .. " changing standing by: " .. mission.standingDifferential .. " for " .. faction)
			local factionObject = query("gameSession","getSessiongOH",faction)[1]
			send(factionObject,
				"changeStanding",
				mission.standingDifferential,
				mission.name)
		end

		function foreign_office_reset_supply_text()
			
			local status = ""
			local supply_warning = ""
			
			if state.supplies[1] >= EntityDB[ state.entityName ].lc_resupply_when_below then
				status = "Working. Office is supplied."
				if state.resupply == false then
					status = "Working. Workcrew ordered to NOT re-supply office."
				end
				SELF.tags.no_supplies1 = nil
				SELF.tags.needs_resupply1 = nil
				SELF.tags.needs_resupply1_badly = nil
				
			elseif state.supplies[1] >= EntityDB[ state.entityName ].mc_resupply_when_below then
				status = "Working. Supplies low. Assigned labourers will re-supply office."
				if state.resupply == false then
					status = "Working. Supplies low. Workcrew ordered to NOT re-supply office."
					SELF.tags.needs_resupply1 = nil
					SELF.tags.needs_resupply1_badly = nil
				else
					SELF.tags.needs_resupply1 = true
					SELF.tags.needs_resupply1_badly = nil
				end
				SELF.tags.no_supplies1 = nil
				supply_warning = "Running low on Bureaucratic Forms."
				
			elseif state.supplies[1] == 0 then
				status = "Work halted. Bureaucratic Forms needed. Workcrew will re-supply office."
				if state.resupply == false then
					status = "Work halted. Workcrew ordered to NOT re-supply office."
					SELF.tags.needs_resupply1 = nil
					SELF.tags.needs_resupply1_badly = nil
				else
					SELF.tags.needs_resupply1 = true
					SELF.tags.needs_resupply1_badly = true
					
					if state.buildingOwner then
						local ownername = query(state.buildingOwner,"getName")[1]
						local alertstring = "The Foreign Office operated by " .. ownername .. " is out of Bureaucratic Forms! Produce more Bureaucratic Forms so diplomatic efforts can continue."
									
						send("rendCommandManager",
							"odinRendererStubMessage", --"odinRendererStubMessage",
							"ui\\orderIcons.xml", -- iconskin
							"foreign_office", -- icon
							"Foreign Office Unsupplied", -- header text
							alertstring, -- text description
							"Left-click to zoom. Right-click to dismiss.", -- action string
							"foreignOfficeProblem", -- alert typasde (for stacking)
							"ui\\eventart\\cult_ritual.png", -- imagename for bg
							"low", -- importance: low / high / critical
							state.rOH, -- object ID
							60 * 1000, -- duration in ms
							0, -- snooze
							state.director)
					end
				end
				
				SELF.tags.no_supplies1 = true
				supply_warning = "Bureaucratic Forms required to do work."
			end
			
			send("rendUIManager", "SetOfficeString", SELF, "noSuppliesWarning",supply_warning)
			send("rendUIManager", "SetOfficeString", SELF, "workPointsStatus", status)
		end
	>>

	state
	<<
		table diplomacyMissions
		table diplomacyPoints
		string currentFaction
		int desksPresent
		int offendedTimes

	>>

	receive Create( stringstringMapHandle init )
	<<
		state.currentFaction = "Empire"
		state.diplomacyPoints.Stahlmark = 0
		state.diplomacyPoints.Novorus = 0
		state.diplomacyPoints.Republique = 0
		state.diplomacyPoints.Empire = 0
		state.diplomacyPoints.Bandits = 0
		state.desksPresent = 0
		state.offendedTimes = 0
		state.diplomacyMissions = {}
		
		-- world->sessionInts
	>>
	
	receive odinBuildingCompleteMessage( int handle, gameSimJobInstanceHandle ji )
	<<
		send(SELF,"generateDiplomaticMissions")
		
		send("rendUIManager", "SetOfficeInt", SELF, "Diplomacy_Points_Fishpeople", 0)
		send("rendUIManager", "SetOfficeInt", SELF, "Diplomacy_Points_Empire"    , 0)
		send("rendUIManager", "SetOfficeInt", SELF, "Diplomacy_Points_Stahlmark" , 0)
		send("rendUIManager", "SetOfficeInt", SELF, "Diplomacy_Points_Novorus"   , 0)
		send("rendUIManager", "SetOfficeInt", SELF, "Diplomacy_Points_Republique", 0)
		send("rendUIManager", "SetOfficeInt", SELF, "Diplomacy_Points_Bandits"   , 0)
		send("rendUIManager", "SetOfficeInt", SELF, "Diplomacy_Points_Fishpeople", 0)

		foreign_office_reset_supply_text()
	>>

	receive InteractiveMessage( string messagereceived )
	<<
		printl ("buildings", "(Foreign Office) Message Received: " .. messagereceived )
		if not state.completed then
			return
		end
		
		if messagereceived == "build" then
			send("rendUIManager", "SetOfficeMenuHeader", SELF, state.buildingName, "overseer_icon")			local counter = 0
		elseif messagereceived == "supply_on" then
			state.resupply = true
			foreign_office_reset_supply_text()
		elseif messagereceived == "supply_off" then
			state.resupply = false
			foreign_office_reset_supply_text()
		elseif messagereceived == "ResetMissions" then
			send (SELF, "generateDiplomaticMissions")
		else
			local factionsUpper = { "Stahlmark", "Novorus", "Empire", "Bandits", "Republique" }
			
			for k,faction in pairs(factionsUpper) do
					
				if string.find(messagereceived,"_mission") and
						string.find(messagereceived, faction) then

					local results1, results2 = string.find(messagereceived, "_mission")
					--local num = string.match(messagereceived, "%d+)
					local missionNumber = tonumber( string.sub(messagereceived, results2 + 1) )
					local faction = string.sub(messagereceived, 1, results1-1)
					
					local factionLC = faction:lower()
					
					
					local missionEntry = state.diplomacyMissions[factionLC][missionNumber]
					
					if missionEntry.enabled then
						
						-- is mission valid ?
						local valid = true
						
						if missionEntry.require_bools_false and missionEntry.require_bools_false[1] then
							valid = false
							for k,v in pairs(missionEntry.require_bools_false) do
								if query("gameSession","getSessionBool",v)[1] == false then
									valid = true
								else
									valid = false
								end
							end
						end
						
						if missionEntry.require_bools_true and missionEntry.require_bools_true[1] then
							valid = false
							for k,v in pairs(missionEntry.require_bools_true) do
								if query("gameSession","getSessionBool",v)[1] == true then
									valid = true
								else
									valid = false
								end
							end
						end
						
						if valid then
							
							if missionEntry.event ~= "" then
								local eventQ = query("gameSimEventManager",
											"startEvent",
											state.diplomacyMissions[factionLC][missionNumber].event,
											{},
											{} )[1]
							elseif missionEntry.standingDifferential then
								
								local faction_info = EntityDB[faction .. "Info"]
								
								local ending_standing = query("gameSession","getSessionInt", faction_info.shortName .. "Relations")[1] + missionEntry.standingDifferential
								
								local alertstring = "Standing with the " .. faction_info.fullName ..
									" has changed by (" .. missionEntry.standingDifferential ..
									") for a result of (" .. ending_standing ..
									") due to the Foreign Office action '" .. missionEntry.name .. "'."
								
								send("rendCommandManager",
									"odinRendererStubMessage", --"odinRendererStubMessage",
									faction_info.iconSkin, -- iconskin
									faction_info.icon, -- icon
									missionEntry.name, -- header text
									alertstring, -- text description
									"Right-click to dismiss.", -- action string
									"foreignOfficeAction", -- alert type (for stacking)
									"ui\\eventart\\capitalists.png", -- imagename for bg
									"low", -- importance: low / high / critical
									state.rOH, -- object ID
									60 * 1000, -- duration in ms
									0, -- snooze
									state.director)
							
							end
					
							missionEntry.enabled = false
							state.diplomacyPoints[faction] = state.diplomacyPoints[faction] - missionEntry.diplomacyCost
							
							send("rendUIManager", "SetOfficeInt", SELF, "Diplomacy_Points_" .. faction, state.diplomacyPoints[faction])
				
							local officeBoolString = factionLC .. "MissionEnabled" .. missionNumber
							send("rendUIManager", "SetOfficeBool", SELF, officeBoolString, false)
				
							changeStanding(faction, missionEntry )
						else
							local alertstring = "The requested Foreign Office action, '" .. missionEntry.name .. "' was unable to be carried out for what you are assured are 'perfectly good reasons'. \z
								Your bureaucrat has saved the forms already written to assist in another Foreign Office action."
								
							send("rendCommandManager",
								"odinRendererStubMessage", --"odinRendererStubMessage",
								"ui\\orderIcons.xml", -- iconskin
								"foreign_office", -- icon
								"Foreign Office action untenable", -- header text
								alertstring, -- text description
								"Right-click to dismiss.", -- action string
								"foreignOfficeProblem", -- alert type (for stacking)
								"ui\\eventart\\cult_ritual.png", -- imagename for bg
								"low", -- importance: low / high / critical
								nil, -- object ID
								60 * 1000, -- duration in ms
								0, -- snooze
								state.director)
						end
					end
				end

				
				-- if message is a faction name, switch office to that faction
				if messagereceived == ("Faction_" .. faction) and state.currentFaction ~= faction then
					eraseDiplomacyPoints() -- and you lose your points, creep.
					state.currentFaction = string.sub(messagereceived,9)
					return
				elseif messagereceived == ("Faction_" .. faction) then
					-- switched to same faction? Alright, there you go.
					return
				end
				
				-- if message is for the selected faction and contains "_mission", do that mission.
			end
		end	
	>>

	receive IncrementDiplomacy( gameObjectHandle bureaucrat, int amount)
	<<
          local AIBlock = {}
          if bureaucrat ~= nil then
               AIBlock = query(bureaucrat, "getAIAttributes")[1]
		end
		if (bureaucrat ~= nil) and (AIBlock.strs.socialClass == "middle") and AIBlock.traits["Xenophobic"] then --uh oh
			
			if state.offendedTimes > 4 then --uh oh!!!
				send("rendCommandManager",
					"odinRendererStubMessage",
					"ui\\thoughtIcons.xml",
					"xenophobic",
					"Diplomatic Incident!", -- header text
					AIBlock.name .. ", the overseer in charge of your foreign office, is causing problems - their Xenophobic \z
					tendancies have caused them to offend nearly every nation we have interactions with! It might be best to reassign them.", -- text description
					"Right-click to dismiss.", -- action string
					"xenophobicDiplomat", -- alert type (for stacking)
					"ui//eventart//bureaucracy_and_fire.png", -- imagename for bg
					"high", -- importance: low / high / critical
					state.renderHandle, -- object ID
					60000, -- duration in ms
					0, -- "snooze" time if triggered multiple times in rapid succession
					nullHandle) -- gameobjecthandle of director, null if none
				
				send( query("gameSession","getSessiongOH", "Stahlmark")[1],
					"changeStanding",-15,nil)
				send( query("gameSession","getSessiongOH", "Republique")[1],
					"changeStanding",-15,nil)
				send( query("gameSession","getSessiongOH", "Novorus")[1],
					"changeStanding",-15,nil)
				
				state.offendedTimes = 0
			else
				state.offendedTimes = state.offendedTimes + 1
			end
			
			state.diplomacyPoints[state.currentFaction] = state.diplomacyPoints[state.currentFaction] - 2
			if state.diplomacyPoints[state.currentFaction] < 0 then
				state.diplomacyPoints[state.currentFaction] = 0
			end
			
		elseif state.diplomacyPoints[state.currentFaction] >= 20 then --You're full on diplomacy!
			
			if bureaucrat then
				local btags = query(bureaucrat, "getTags")[1]
				if btags.middle_class then 
					local name = query(bureaucrat,"getName")[1]
					
					local tickerText = name.. " 's Foreign Office is completely packed with Bureaucracy and can't fit any more! \z
						Spend some Diplomacy points, or " .. name .. " will just keep doing the Empire Times crossword puzzle \z
						instead of doing useful work."
					
					send("rendCommandManager",
						"odinRendererFYIMessage",
						"ui\\orderIcons.xml",
						"foreign_office",
						"Diplomacy Points Full", -- header text
						tickerText, -- text description
						"Left-click for details. Right-click to dismiss.", -- action string
						"diplomacyFull", -- alert type (for stacking)
						"ui//eventart//capitalists.png", -- imagename for bg
						"high", -- importance: low / high / critical
						state.rOH, -- object ID
						30 * 1000, -- duration in ms
						0, -- "snooze" time if triggered multiple times in rapid succession
						nil) -- gameobjecthandle of director, null if none
				end
			end
		else
			if (bureaucrat ~= nil) and (state.currentFaction == "Empire") and
				AIBlock.strs.socialClass == "middle" and
				AIBlock.traits["Patriotic"] then
				--Huzzah!!
				
				state.diplomacyPoints[state.currentFaction] = state.diplomacyPoints[state.currentFaction] + (amount * 2)
			else
				--This is the normal result. Everything else is... special cases. Yay!
				state.diplomacyPoints[state.currentFaction] = state.diplomacyPoints[state.currentFaction] + amount
			end
		end
		
		send("rendUIManager",
			"SetOfficeInt",
			SELF,
			"Diplomacy_Points_" .. state.currentFaction,
			state.diplomacyPoints[state.currentFaction])
		
	>>

	receive newShiftUpdate()
	<<
		local currentShift = query("gameSession", "getSessionInt", "currentShift")[1]
		if currentShift == 1 then
			send(SELF, "generateDiplomaticMissions")
		end
	>>

	receive generateDiplomaticMissions()
	<<
		local factionsLC = { "empire", "bandits", "stahlmark", "novorus", "republique"}
		
		for k,faction in pairs(factionsLC) do
			
			local factionUpper = faction:gsub("^%l", string.upper)
			local factionInfo = EntitiesByType["faction"][factionUpper .. "Info"]
			local randoms = {}
			state.diplomacyMissions[faction] = {}
			
			printl("buildings", "foreign office rebuilding missions for: " .. tostring(factionInfo.shortName) )
			
			-- only choose missions that work w/ current relations state
			-- ie. if "friendly" only choose "friendly" missions
			local missionSet = {}
			
			if factionInfo.shortName == "Novorus" or
				factionInfo.shortName == "Republique" or
				factionInfo.shortName == "Stahlmark" or
				factionInfo.shortName == "Bandits" then
			
				-- use standing to determine missions set
				local factionState = query( query("gameSession",
										    "getSessiongOH",
										    factionInfo.shortName)[1],
											"getRelationStateString")[1]
				
				--printl("DAVID", "faction state is: " .. tostring(factionState) )
				
				for k,v in pairs(factionInfo.missions) do
					if v.standing == factionState or
						v.standing == "any" then
						missionSet[ #missionSet+1 ] = v 
					end
				end
				
			else
				-- yolo it
				-- but TODO: make this work correctly for other faction types.

				for k,v in pairs(factionInfo.missions) do
					missionSet[ #missionSet+1 ] = v 
				end
			end
			
			if #missionSet < 3 then
				printl("office", "WARNING: low missionset!!! just ... use neutral missions?")
				missionSet = factionInfo.missions
			end
			
			local i = 1
			local backup_count = 0
			while i ~= 4 do 
				state.diplomacyMissions[faction][i] = {}
				
				local valid = true
				local randomNum = rand(1,#missionSet)
				local missionEntry = table.remove(missionSet, randomNum)
				
				if missionEntry.require_bools_false and missionEntry.require_bools_false[1] then
					valid = false
					for k,v in pairs(missionEntry.require_bools_false) do
						if query("gameSession","getSessionBool",v)[1] == false then
							valid = true
						else
							valid = false
						end
					end
				end
				
				if missionEntry.require_bools_true and missionEntry.require_bools_true[1] then
					valid = false
					for k,v in pairs(missionEntry.require_bools_true) do
						if query("gameSession","getSessionBool",v)[1] == true then
							valid = true
						else
							valid = false
						end
					end
				end
				
				if valid == true then
					for k, v in pairs( missionEntry ) do
						state.diplomacyMissions[faction][i][k] = v
					end
					
					state.diplomacyMissions[faction][i].enabled = true
					i = i + 1
					randoms[i] = randomNum
					backup_count = 0
				end
				
				if backup_count > 10 then
					randoms = {}
				end
			end
	
			for count=1,3 do
				send("rendUIManager", "SetOfficeInt", SELF, faction .. "MissionCost" .. count, state.diplomacyMissions[faction][count].diplomacyCost)
				send("rendUIManager", "SetOfficeString", SELF, faction .. "MissionName" .. count, state.diplomacyMissions[faction][count].name)
				send("rendUIManager", "SetOfficeString", SELF, faction .. "MissionDesc" .. count, state.diplomacyMissions[faction][count].description)
				send("rendUIManager", "SetOfficeBool", SELF, faction .. "MissionEnabled" .. count, state.diplomacyMissions[faction][count].enabled)
				send("rendUIManager", "SetOfficeInt", SELF, faction .. "MissionStandingChange" .. count, state.diplomacyMissions[faction][count].standingDifferential)
				send("rendUIManager", "SetOfficeString", SELF, faction .. "MissionIconSkin" .. count, state.diplomacyMissions[faction][count].iconSkin)
				send("rendUIManager", "SetOfficeString", SELF, faction .. "MissionIcon" .. count, state.diplomacyMissions[faction][count].icon)
			end
		end
	>>
	
	receive BurnDiplomaticPapers()
	<<
		printl("events", "foreign office received order to burn diplomatic papers!")
		
		state.diplomacyPoints[state.currentFaction] = 0
		send("rendUIManager", "SetOfficeInt", SELF, "Diplomacy_Points_" .. state.currentFaction, state.diplomacyPoints[state.currentFaction])
		
		--[[send(SELF, "generateDiplomaticMissions")
		
		local factionsLC = { "empire", "bandits", "stahlmark", "novorus", "republique"}
		for k,faction in pairs(factionsLC) do
			for count=1,3 do
				send("rendUIManager", "SetOfficeBool", SELF, faction .. "MissionEnabled" .. count, false) -- disable all missions today.
			end
		end]]
	>>
	
	receive buildingAddSupplies( gameObjectHandle item, int count)
	<<
		local tags = query(item,"getTags")[1]
		local mult = 1
		local tier = 1

		for k,v in pairs( EntityDB[ state.entityName ].supply_info ) do
			if tags[k] then
				mult = v.multiplier
				tier = v.tier
				break
			end
		end

		state.supplies[tier] = state.supplies[tier] + count * mult

		send("rendUIManager", "SetOfficeInt", SELF, "workPoints1", state.supplies[1])
		foreign_office_reset_supply_text()
	>>
	
	receive consumeSupplies( int tier, int count )
	<<
		state.supplies[tier] = state.supplies[tier] - count 
		send("rendUIManager", "SetOfficeInt", SELF, "workPoints1", state.supplies[tier])
		foreign_office_reset_supply_text()
	>>
	
	receive recalculateQuality()
	<<
		-- count number of desks
		local modules = query("gameWorkshopManager", "getBuildingModulesGameSide", SELF)[1]
		local desk_count = 0
		for k,v in pairs(modules) do
			local module_tags = query(v,"getTags")[1]
			if module_tags.diplomatic_desk then
				desk_count = desk_count + 1
			end
		end
		
		send("rendUIManager", "SetOfficeInt", SELF, "desksPresent", desk_count)
		if desk_count == 0 then
			send("rendUIManager","SetOfficeString",SELF,"noDeskWarning","At least one desk is required to perform work.")
		else
			send("rendUIManager","SetOfficeString",SELF,"noDeskWarning","")
		end
	>>
>>