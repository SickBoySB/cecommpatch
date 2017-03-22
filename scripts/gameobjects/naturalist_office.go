gameobject "naturalist_office" inherit "office"
<<
	local 
	<<
		function naturalist_office_reset_supply_text()
			
			local status = ""
			local supply_warning = ""
			local office_data = EntityDB[ state.entityName ]
			
			if state.mode_scout then
				status = "No supplies required to do Scouting."
			end
			
			if state.supplies[1] >= office_data.lc_resupply_when_below then
				if state.mode_hunt then
					status = "Working. Office is supplied with ammo."
					if state.resupply == false then
						status = "Working. Overseer/Labourers ordered to NOT re-supply office."
					end
				end
				SELF.tags.no_supplies1 = nil
				SELF.tags.needs_resupply1 = nil
				SELF.tags.needs_resupply1_badly = nil
				
			elseif state.supplies[1] > office_data.mc_resupply_when_below then
				if state.mode_hunt then
					status = "Working. Supplies low. Assigned labourers will re-supply office."
					supply_warning = supply_warning .. "Low on hunting ammo."
					if state.resupply == false then
						status = "Working. Overseer/Labourers ordered to NOT re-supply office."
					end
				end
				
				if state.resupply == false then
					SELF.tags.needs_resupply1 = nil
					SELF.tags.needs_resupply1_badly = nil
				else
					SELF.tags.needs_resupply1 = true
					SELF.tags.needs_resupply1_badly = nil
				end
				
				SELF.tags.no_supplies1 = nil
				
			elseif state.supplies[1] == 0 then
				if state.mode_hunt then
					status = "Work halted. Stone Pellet Ammunition needed. Assigned Overseer/Labourers will re-supply office."
					supply_warning = "Stone Pellet Ammunition required to do Hunting."
					if state.resupply == false then
						status = "Work halted. Overseer/Labourers ordered to NOT re-supply office."
					else
						if state.buildingOwner then
							local ownername = query(state.buildingOwner,"getName")[1]
							local alertstring = "The Naturalist's Office operated by " .. ownername .. " is out of Stone Pellet Ammunition! Produce more to enable hunting jobs."
							
							send("rendCommandManager",
								"odinRendererStubMessage", --"odinRendererStubMessage",
								"ui\\orderIcons.xml", -- iconskin
								"hunting", -- icon
								"Naturalist Needs Ammo", -- header text
								alertstring, -- text description
								"Left-click to zoom. Right-click to dismiss.", -- action string
								"naturalistOfficeProblem", -- alert typasde (for stacking)
								"ui\\eventart\\cult_ritual.png", -- imagename for bg
								"low", -- importance: low / high / critical
								state.rOH, -- object ID
								60 * 1000, -- duration in ms
								0, -- snooze
								state.director)
						end
					end
				end
				
				if state.resupply == false then
					SELF.tags.needs_resupply1 = nil
					SELF.tags.needs_resupply1_badly = nil
				else
					SELF.tags.needs_resupply1 = true
					SELF.tags.needs_resupply1_badly = true
				end
				
				SELF.tags.no_supplies1 = true
			end
			
			if state.supplies[2] >= office_data.lc_resupply_when_below then
				if state.mode_survey then
					status = "Working. Office is supplied with paperwork."
					if state.resupply == false then
						status = "Working.  Overseer/Labourers ordered to NOT re-supply office."
					end
				end
				SELF.tags.no_supplies2 = nil
				SELF.tags.needs_resupply2 = nil
				SELF.tags.needs_resupply2_badly = nil
				
			elseif state.supplies[2] > office_data.mc_resupply_when_below then
				if state.mode_survey then
					status = "Working. Supplies low. Assigned labourers will re-supply office."
					supply_warning = supply_warning .. "Low on Bureaucratic Forms."
					if state.resupply == false then
						status = "Working. Overseer/Labourers ordered to NOT re-supply office."
					end
				end
				
				if state.resupply == false then
					SELF.tags.needs_resupply2 = nil
					SELF.tags.needs_resupply2_badly = nil
				else
					SELF.tags.needs_resupply2 = true
					SELF.tags.needs_resupply2_badly = nil
				end
				
				SELF.tags.no_supplies2 = nil
				
			elseif state.supplies[2] == 0 then
				
				if state.mode_survey then
					status = "Work halted. Bureaucratic Forms needed. Overseer/Labourers will re-supply office."
					supply_warning = "Bureaucratic Forms required to do Surveying."
					if state.resupply == false then
						status = "Work halted. Overseer/Labourers ordered to NOT re-supply office."
					else
						if state.buildingOwner then
							local ownername = query(state.buildingOwner,"getName")[1]
							local alertstring = "The Naturalist's Office operated by " .. ownername .. " is out of Stone Pellet Ammunition! Produce more to enable hunting jobs."
							
							send("rendCommandManager",
								"odinRendererStubMessage", --"odinRendererStubMessage",
								"ui\\commodityIcons.xml", -- iconskin
								"paperBundle", -- icon
								"Naturalist Needs Forms", -- header text
								alertstring, -- text description
								"Left-click to zoom. Right-click to dismiss.", -- action string
								"naturalistOfficeProblem", -- alert typasde (for stacking)
								"ui\\eventart\\cult_ritual.png", -- imagename for bg
								"low", -- importance: low / high / critical
								state.rOH, -- object ID
								60 * 1000, -- duration in ms
								0, -- snooze
								state.director)
						end
					end
				end
				
				if state.resupply == false then
					SELF.tags.needs_resupply2 = nil
					SELF.tags.needs_resupply2_badly = nil
				else
					SELF.tags.needs_resupply2 = true
					SELF.tags.needs_resupply2_badly = true
				end
				
				SELF.tags.no_supplies2 = true
			end
			
			send("rendUIManager", "SetOfficeString", SELF, "noSuppliesWarning",supply_warning)
			send("rendUIManager", "SetOfficeString", SELF, "workPointsStatus", status)
		end
	>>

	state
	<<
		string currentTask
		int targetMapCell
	>>

	receive Create( stringstringMapHandle init )
	<<
		send("rendUIManager", "SetOfficeString", SELF, "modeLabel", "Scouting")
	
		state.currentTask = "Scout"
		state.targetMapCell = 6
	>>
	
	receive odinBuildingCompleteMessage( int handle, gameSimJobInstanceHandle ji )
	<<
		SELF.tags.mode_scout = true
		SELF.tags.mode_survey = nil
		SELF.tags.mode_hunt = nil
		SELF.tags.mode_discover = nil
			
		send("rendUIManager", "SetOfficeInt", SELF, "workPoints1", 0)
		send("rendUIManager", "SetOfficeInt", SELF, "workPoints2", 0)
		send("rendUIManager", "SetOfficeString", SELF, "modeLabel", "Scouting")
		naturalist_office_reset_supply_text()
		
		send(SELF,"InteractiveMessage","mode_scout")
		
		if query("gameSession","getSessionBool","horror_policy_study")[1] then
			
			send("rendUIManager", "SetOfficeString", SELF, "horrorPolicyLabel", "Study")
			
		elseif query("gameSession","getSessionBool","horror_policy_harvest")[1] then
			
			send("rendUIManager", "SetOfficeString", SELF, "horrorPolicyLabel", "Harvest")
			
		elseif query("gameSession","getSessionBool","horror_policy_dump")[1] then
			
			send("rendUIManager", "SetOfficeString", SELF, "horrorPolicyLabel", "Destroy")
		else
			
			send("rendUIManager", "SetOfficeString", SELF, "horrorPolicyLabel", "Destroy")
		end

	>>
	
	receive InteractiveMessage( string messagereceived )
	<<
		printl("buildings", "naturalist office received message: " .. messagereceived)
		if not state.completed or SELF.tags.slated_for_demolition then
			return
		end
		
		if messagereceived == "mode_scout" then
			SELF.tags.mode_scout = true
			SELF.tags.mode_survey = nil
			SELF.tags.mode_hunt = nil
			SELF.tags.mode_discover = nil
			
			send("rendUIManager", "SetOfficeString", SELF, "modeLabel", "Scouting")
			naturalist_office_reset_supply_text()
			
		elseif messagereceived == "mode_survey" then
			SELF.tags.mode_scout = nil
			SELF.tags.mode_survey = true
			SELF.tags.mode_hunt = nil
			SELF.tags.mode_discover = nil
			
			send("rendUIManager", "SetOfficeString", SELF, "modeLabel", "Surveying for Minerals")
			naturalist_office_reset_supply_text()
			
			if not state.buildingOwner then
				local alertstring = "A Naturalist's Office has been ordered to beging Surveying for Minerals, but it has no assigned Overseer. Assign an Overseer to begin work."
				
				send("rendCommandManager",
					"odinRendererStubMessage", --"odinRendererStubMessage",
					"ui\\orderIcons.xml", -- iconskin
					"naturalism_icon", -- icon
					"Office needs Overseer!", -- header text
					alertstring, -- text description
					"Left-click to zoom. Right-click to dismiss.", -- action string
					"naturalistProblem", -- alert typasde (for stacking)
					"ui\\eventart\\naturalistAndButterflies.png", -- imagename for bg
					"low", -- importance: low / high / critical
					state.rOH, -- object ID
					60 * 1000, -- duration in ms
					0, -- snooze
					state.director)
				
			elseif state.buildingOwner and state.supplies[2] == 0 then
				
				local ownername = query(state.buildingOwner,"getName")[1]
				local alertstring = "The Naturalist's Office operated by " .. ownername .. " requires Paper Bundles! Produce some Paper Bundles in a Carpentry Workshop so surveying may proceed."
			
				send("rendCommandManager",
					"odinRendererStubMessage", --"odinRendererStubMessage",
					"ui\\commodityIcons.xml", -- iconskin
					"paperBundle", -- icon
					"Naturalist Needs Paper", -- header text
					alertstring, -- text description
					"Left-click to zoom. Right-click to dismiss.", -- action string
					"naturalistProblem", -- alert typasde (for stacking)
					"ui\\eventart\\naturalistAndButterflies.png", -- imagename for bg
					"low", -- importance: low / high / critical
					state.rOH, -- object ID
					60 * 1000, -- duration in ms
					0, -- snooze
					state.director)
			end
			
		elseif messagereceived == "mode_hunt" then
			SELF.tags.mode_scout = nil
			SELF.tags.mode_survey = nil
			SELF.tags.mode_hunt = true
			SELF.tags.mode_discover = nil
			 
			send("rendUIManager", "SetOfficeString", SELF, "modeLabel", "Hunting")
			naturalist_office_reset_supply_text()
			
			if not state.buildingOwner then
				local alertstring = "A Naturalist's Office has been ordered to beging Hunting but it has no assigned Overseer. Assign an Overseer to begin work."
				
				send("rendCommandManager",
					"odinRendererStubMessage", --"odinRendererStubMessage",
					"ui\\orderIcons.xml", -- iconskin
					"naturalism_icon", -- icon
					"Office needs Overseer!", -- header text
					alertstring, -- text description
					"Left-click to zoom. Right-click to dismiss.", -- action string
					"naturalistProblem", -- alert typasde (for stacking)
					"ui\\eventart\\naturalistAndButterflies.png", -- imagename for bg
					"low", -- importance: low / high / critical
					state.rOH, -- object ID
					60 * 1000, -- duration in ms
					0, -- snooze
					state.director)
				
			elseif state.buildingOwner and state.supplies[1] == 0 then
				
				local ownername = query(state.buildingOwner,"getName")[1]
				local alertstring = "The Naturalist's Office operated by " .. ownername .. " requires Stone Pellet Ammunition! Produce this ammo in a Ceramics Workshop so hunting may proceed."
			
				send("rendCommandManager",
					"odinRendererStubMessage", --"odinRendererStubMessage",
					"ui\\commodityIcons.xml", -- iconskin
					"ammo_stone", -- icon
					"Nautralist Needs Ammo", -- header text
					alertstring, -- text description
					"Left-click to zoom. Right-click to dismiss.", -- action string
					"naturalistProblem", -- alert typasde (for stacking)
					"ui\\eventart\\naturalistAndButterflies.png", -- imagename for bg
					"low", -- importance: low / high / critical
					state.rOH, -- object ID
					60 * 1000, -- duration in ms
					0, -- snooze
					state.director)
			end
			
		elseif messagereceived == "mode_discovery" then
			SELF.tags.mode_scout = nil
			SELF.tags.mode_survey = nil
			SELF.tags.mode_hunt = nil
			SELF.tags.mode_discover = true
			
			send("rendUIManager", "SetOfficeString", SELF, "modeLabel", "Performing Naturalism Research")
			naturalist_office_reset_supply_text()
		elseif messagereceived == "supply_on" then
			state.resupply = true
			naturalist_office_reset_supply_text()
		elseif messagereceived == "supply_off" then
			state.resupply = false
			naturalist_office_reset_supply_text()
			
		elseif messagereceived == "horror_policy_study" then
			printl("buildings", state.buildingName .. " got horror_policy_study" )
			-- send this to all naturalist offices
			local results = query("gameSpatialDictionary", "allBuildingsRequest")
			if results and results[1] then
				for k,v in pairs(results[1]) do
					send(v, "setHorrorPolicy","study")
				end
			end
			
			send("gameSession","setSessionBool","horror_policy_study", true)
			send("gameSession","setSessionBool","horror_policy_harvest", false)
			send("gameSession","setSessionBool","horror_policy_dump", false)
			
		elseif messagereceived == "horror_policy_harvest" then
			printl("buildings", state.buildingName .. " got horror_policy_harvest" )
			-- send this to all naturalist offices
			local results = query("gameSpatialDictionary", "allBuildingsRequest")
			if results and results[1] then
				for k,v in pairs(results[1]) do
					send(v, "setHorrorPolicy","harvest")
				end
			end
			
			send("gameSession","setSessionBool","horror_policy_study", false)
			send("gameSession","setSessionBool","horror_policy_harvest", true)
			send("gameSession","setSessionBool","horror_policy_dump", false)
			
		elseif messagereceived == "horror_policy_dump" then
			printl("buildings", state.buildingName .. " got horror_policy_dump" )
			-- send this to all naturalist offices
			local results = query("gameSpatialDictionary", "allBuildingsRequest")
			if results and results[1] then
				for k,v in pairs(results[1]) do
					send(v, "setHorrorPolicy","dump")
				end
			end
			
			send("gameSession","setSessionBool","horror_policy_study", false)
			send("gameSession","setSessionBool","horror_policy_harvest", false)
			send("gameSession","setSessionBool","horror_policy_dump", true)
			
		end
	>>
	
	receive setWeaponLoadout(string loadoutname)
	<<
		printl("buildings", "barracks got setWeaponLoadout: " .. tostring(loadoutname) )
		if loadoutname then
			state.weaponLoadout = loadoutname
		else
			state.weaponLoadout = "pistol"
		end
	>>
	
	respond getWeaponLoadout()
	<<
		return "weaponLoadoutMessage", state.weaponLoadout
	>>
	
	receive setBuildingOwner(gameObjectHandle newOwner)
	<<
		-- set new owner to correct shift.
		--[[if newOwner then
			if not state.currentShiftSelection then state.currentShiftSelection = 1 end
			--send(newOwner,"InteractiveMessage","setStartHour" .. tostring(state.currentShiftSelection))
			--send("rendCommandManager", "gameSetWorkPartyWorkShift", newOwner, state.currentShiftSelection)
			
			if state.currentShiftSelection == 1 then
				send(newOwner,"setWorkShift", 0, 0, 0, 0, 2, 5, 1, 1) 
			elseif state.currentShiftSelection == 5 then
				send(newOwner,"setWorkShift", 2, 5, 1, 1, 0, 0, 0, 0)
			end
		end]]
	>>

	receive recalculateQuality()
	<<
		local desk = false
		local muskets = false
		
		local modules = query("gameWorkshopManager", "getBuildingModulesGameSide", SELF)[1]
		for k,v in pairs(modules) do
			for tagname,tagbool in pairs(v.tags) do
				if tagname == "standing_desk" then
					desk = true
				elseif tagname == "musket_locker" then
					muskets = true
				end
			end
		end
		
		if muskets then
			send("rendUIManager", "SetOfficeBool", SELF, "huntEnabled", true)
			state.has_muskets = true
		else
			if state.mode == "mode_hunt" then
				send(SELF,"InteractiveMessage","mode_scout")
				
				local alertstring = "Your Naturalist had orders to hunt but the Musket Locker was removed! Orders have been reverted to scouting."
				send("rendCommandManager",
					"odinRendererStubMessage", --"odinRendererStubMessage",
					"ui\\orderIcons.xml", -- iconskin
					"naturalism_icon", -- icon
					"Can't Hunt!", -- header text
					alertstring, -- text description
					"Left-click to zoom. Right-click to dismiss.", -- action string
					"naturalistOfficeProblem", -- alert typasde (for stacking)
					"ui\\eventart\\cult_ritual.png", -- imagename for bg
					"low", -- importance: low / high / critical
					state.rOH, -- object ID
					60 * 1000, -- duration in ms
					0, -- snooze
					state.director)
			end
			send("rendUIManager", "SetOfficeBool", SELF, "huntEnabled", false)
			state.has_muskets = false
		end
		
		if desk then
			send("rendUIManager", "SetOfficeBool", SELF, "surveyEnabled", true)
			state.has_desk = true
		else
			if state.mode == "mode_survey" then
				send(SELF,"InteractiveMessage","mode_scout")
				
				local alertstring = "Your Naturalist had orders to do mineralogical surveyis but the Standing Desk was removed! Orders have been reverted to scouting."
				send("rendCommandManager",
					"odinRendererStubMessage", --"odinRendererStubMessage",
					"ui\\orderIcons.xml", -- iconskin
					"naturalism_icon", -- icon
					"Can't Survey!", -- header text
					alertstring, -- text description
					"Left-click to zoom. Right-click to dismiss.", -- action string
					"naturalistOfficeProblem", -- alert typasde (for stacking)
					"ui\\eventart\\cult_ritual.png", -- imagename for bg
					"low", -- importance: low / high / critical
					state.rOH, -- object ID
					60 * 1000, -- duration in ms
					0, -- snooze
					state.director)
			end
			send("rendUIManager", "SetOfficeBool", SELF, "surveyEnabled", false)
			state.has_desk = false
		end
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
		
		send("rendUIManager", "SetOfficeInt", SELF, "workPoints" .. tostring(tier), state.supplies[tier])
		naturalist_office_reset_supply_text()
	>>
	
	receive consumeSupplies( int tier, int count )
	<<
		state.supplies[tier] = state.supplies[tier] - count 
		send("rendUIManager", "SetOfficeInt", SELF, "workPoints" .. tostring(tier), state.supplies[tier])
		naturalist_office_reset_supply_text()
	>>
	
	receive setHorrorPolicy( string policy )
	<<
		if policy == "study" then
			send("rendUIManager", "SetOfficeString", SELF, "horrorPolicyLabel", "Study")
		elseif policy == "harvest" then
			send("rendUIManager", "SetOfficeString", SELF, "horrorPolicyLabel", "Harvest")
		else
			send("rendUIManager", "SetOfficeString", SELF, "horrorPolicyLabel", "Destroy")
		end
	>>
>>