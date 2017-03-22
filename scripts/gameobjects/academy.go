gameobject "academy" inherit "office"
<<
	local 
	<<
		function academy_reset_supply_text()
			
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
						status = "Working. Workcrew ordered to NOT re-supply office."
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
						status = "Working. Workcrew ordered to NOT re-supply office."
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
						status = "Work halted. Workcrew ordered to NOT re-supply office."
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
						status = "Working. Workcrew ordered to NOT re-supply office."
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
						status = "Working. Workcrew ordered to NOT re-supply office."
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
					status = "Work halted. Bureaucratic Forms needed. Assigned Overseer/Labourers will re-supply office."
					supply_warning = "Bureaucratic Forms required to do Surveying."
					if state.resupply == false then
						status = "Work halted. Workcrew ordered to NOT re-supply office."
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
		string currentSkill
		int bookShelvesPresent
	>>

	receive Create( stringstringMapHandle init )
	<<
		state.currentSkill = "carpentry"
	>>
	
	receive odinBuildingCompleteMessage ( int handle, gameSimJobInstanceHandle ji )
	<<
		send("rendUIManager", "SetOfficeInt", SELF, "workPoints1", 0)
		send("rendUIManager", "SetOfficeInt", SELF, "workPoints2", 0)
		send("rendUIManager","SetOfficeInt",SELF,"bookshelvesPresent",0)
		send("rendUIManager","SetOfficeString",SELF,"noBookshelvesWarning","At least one bookshelf is required to train.")
		academy_reset_supply_text()
	>>
	
	receive InteractiveMessage( string messagereceived )
	<<
		if not state.completed then
			return
		end
		if messagereceived == "Skill_button1" then
			state.currentSkill = "carpentry"
		elseif messagereceived == "Skill_button2" then
			state.currentSkill = "naturalism"
		elseif messagereceived == "Skill_button3" then
			state.currentSkill = "smithing"
		elseif messagereceived == "Skill_button4" then
               state.currentSkill = "stoneworking"
          elseif messagereceived == "Skill_button5" then
               state.currentSkill = "cooking"
          elseif messagereceived == "Skill_button6" then
               state.currentSkill = "science"
          elseif messagereceived == "Skill_button7" then
               state.currentSkill = "farming"
			
		elseif messagereceived == "supply_on" then
			state.resupply = true
			academy_reset_supply_text()
		elseif messagereceived == "supply_off" then
			state.resupply = false
			academy_reset_supply_text()
		end
	>>
	
	respond GetSkill()
	<<
		return "skillResult", state.currentSkill
	>>
	
	receive recalculateQuality()
	<<
		local shelf_count = 0
          local modules = query("gameWorkshopManager", "getBuildingModulesGameSide", SELF)[1]
		for k,v in pairs(modules) do
			local m_tags = query(v,"getTags")[1]
			if m_tags.bookshelf then
				shelf_count = shelf_count + 1
			end
		end
		
		send("rendUIManager","SetOfficeInt",SELF,"bookshelvesPresent",shelf_count)
		
		if shelf_count > 0 then 
			send("rendUIManager","SetOfficeString",SELF,"noBookshelvesWarning","")
		else
			send("rendUIManager","SetOfficeString",SELF,"noBookshelvesWarning","At least one bookshelf is required to train.")
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

		send("rendUIManager", "SetOfficeInt", SELF, "workPoints" .. tostring(tier) , state.supplies[tier])
		academy_reset_supply_text()
	>>
	
	receive consumeSupplies( int tier, int count )
	<<
		state.supplies[tier] = state.supplies[tier] - count 
		send("rendUIManager", "SetOfficeInt", SELF, "workPoints1", state.supplies[tier])
		academy_reset_supply_text()
	>>
>>