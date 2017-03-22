gameobject "chapel" inherit "office"
<<
	local 
	<<
		function chapel_reset_supply_text()
			
			local status = ""
			local supply_warning = ""
			local office_data = EntityDB[ state.entityName ]
			
			--if state.mode_scout then
				--status = "No supplies required to perform Chapel rituals."
			--end
			
			if state.supplies[1] >= office_data.lc_resupply_when_below then
				--if state.mode_hunt then
					--status = "Working. Chapel is supplied with iron votive cogs."
					if state.resupply == false then
						--status = "Working. Workcrew ordered to NOT re-supply chapel."
					end
				--end
				SELF.tags.no_supplies1 = nil
				SELF.tags.needs_resupply1 = nil
				SELF.tags.needs_resupply1_badly = nil
				
			elseif state.supplies[1] > office_data.mc_resupply_when_below then
				--if state.mode_hunt then
					--status = "Working. Supplies low. Assigned labourers will re-supply chapel."
					--supply_warning = supply_warning .. "Low on iron votive cogs."
					if state.resupply == false then
						--status = "Working. Workcrew ordered to NOT re-supply chapel."
					end
				--end
				
				if state.resupply == false then
					SELF.tags.needs_resupply1 = nil
					SELF.tags.needs_resupply1_badly = nil
				else
					SELF.tags.needs_resupply1 = true
					SELF.tags.needs_resupply1_badly = nil
				end

				SELF.tags.no_supplies1 = nil
				
			elseif state.supplies[1] == 0 then
				--if state.mode_hunt then
					--status = "Work halted. Iron Cogs needed. Assigned Overseer/Labourers will re-supply chapel."
					--supply_warning = "Iron Cogs required to do perform Chapel services."
					status = "Iron Cogs needed!"
					if state.resupply == false then
						--status = "Work halted. Overseer/Labourers ordered to NOT re-supply chapel."
					else
						if state.buildingOwner then
							local ownername = query(state.buildingOwner,"getName")[1]
							local alertstring = "The Chapel operated by " .. ownername .. " is out of Iron Cogs! Produce more to enable Chapel services."
							
							send("rendCommandManager",
								"odinRendererStubMessage", --"odinRendererStubMessage",
								"ui\\orderIcons.xml", -- iconskin
								"chapel", -- icon
								"Chapel needs Cogs", -- header text
								alertstring, -- text description
								"Left-click to zoom. Right-click to dismiss.", -- action string
								"chapelProblem", -- alert typasde (for stacking)
								"ui\\eventart\\cult_ritual.png", -- imagename for bg
								"low", -- importance: low / high / critical
								state.rOH, -- object ID
								60 * 1000, -- duration in ms
								0, -- snooze
								state.director)
						end
					end
				--end
				
				if state.resupply == false then
					SELF.tags.needs_resupply1 = nil
					SELF.tags.needs_resupply1_badly = nil
				else
					SELF.tags.needs_resupply1 = true
					SELF.tags.needs_resupply1_badly = true
				end
				
				SELF.tags.no_supplies1 = true
			end
			
			send("rendUIManager", "SetOfficeString", SELF, "noSuppliesWarning",supply_warning)
			send("rendUIManager", "SetOfficeString", SELF, "workPointsStatus", status)
		end
	>>

	state
	<<
		string currentDoctrine
		int altarsPresent
	>>

	receive Create( stringstringMapHandle init )
	<<
		state.currentDoctrine = "Cog"
	>>
	
	receive InteractiveMessage( string messagereceived )
	<<
		if not state.completed then
			return
		end
		if messagereceived == "Doctrine_button1" then
			state.currentDoctrine = "Cog"
		elseif messagereceived == "Doctrine_button2" then
			state.currentDoctrine = "Military"
		elseif messagereceived == "Doctrine_button3" then
			state.currentDoctrine = "Worker"
		elseif messagereceived == "Doctrine_button4" then
			state.currentDoctrine = "Upperclass"
		end
	>>
	
	respond GetDoctrine()
	<<
		return "doctrineResult", state.currentDoctrine
	>>
	
	receive shiftToggle(string status)
	<<
		if status == "off" then
			send("rendUIManager", "SetOfficeInt", SELF, "shiftStatus", "OFF")
		elseif status == "on" then
			send("rendUIManager", "SetOfficeInt", SELF, "shiftStatus", "ON")
		else
			printl("buildings", "chapel: Office received a bad shift toggle!!")
		end
	>>
	
	receive odinBuildingCompleteMessage ( int handle, gameSimJobInstanceHandle ji )
	<<
		send("rendUIManager","SetOfficeInt",SELF,"workPoints1",0)
		--send("rendUIManager","SetOfficeInt",SELF,"workPoints2",0)
		
		send("rendUIManager","SetOfficeString",SELF,"noAltarWarning","Altar needed!")
		send("rendUIManager","SetOfficeString",SELF,"workPointsStatus","Iron Cogs needed!")
		
		chapel_reset_supply_text()
		
		if not query("gameSession","getSessionBool","builtAChapel")[1] then
			send("gameSession","setSessionBool","builtAChapel",true)
			send("gameSession", "setSteamAchievement", "builtAChapel")
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
		chapel_reset_supply_text()
	>>
	
	receive consumeSupplies( int tier, int count )
	<<
		state.supplies[tier] = state.supplies[tier] - count 
		send("rendUIManager", "SetOfficeInt", SELF, "workPoints" .. tostring(tier), state.supplies[tier])
		chapel_reset_supply_text()
	>>
	
	receive recalculateQuality()
	<<
		local has_chair = false
		local has_altar = false
		local altarCount = 0
          local modules = query("gameWorkshopManager", "getBuildingModulesGameSide", SELF)[1]
		for k,v in pairs(modules) do --TODO: Check for modules that haven't had their upkeep paid
               local module_tags = query(v, "getTags")[1]
               if module_tags.sittable then
				has_chair = true
			elseif module_tags.chapel_altar then
				has_altar = true
				altarCount = altarCount + 1
			end
		end
		
		if has_altar then
			send("rendUIManager","SetOfficeString",SELF,"noAltarWarning","")
		else
			send("rendUIManager","SetOfficeString",SELF,"noAltarWarning","Altar needed!")
		end
		
		send("rendUIManager", "SetOfficeInt", SELF, "altarsPresent",altarCount)	
     >>
	
	receive setBuildingOwner( gameObjectHandle newOwner )
	<<
		if newOwner and state.supplies[1] == 0 then
			
			local eventQ = query("gameSimEventManager",
								"startEvent",
								"supplies_warning_chapel",
								{},
								{} )[1]
						
			send(eventQ,"registerBuilding",SELF)
		end
	>>
>>