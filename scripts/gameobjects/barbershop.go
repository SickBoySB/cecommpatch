gameobject "barbershop" inherit "office"
<<
	local 
	<<
		function barbershop_reset_supply_text()
			
			local status = ""
			local supply_warning = ""
			local office_data = EntityDB[ state.entityName ]
			local status_data = {}
			
			if state.supplies[1] >= office_data.lc_resupply_when_below then
				--status = "Working. Medical Supplies stocked."
				if state.resupply == false then
					--status = "Working. Workcrew ordered to NOT re-supply office."
				end

				SELF.tags.no_supplies1 = nil
				SELF.tags.needs_resupply1 = nil
				SELF.tags.needs_resupply1_badly = nil
				
			elseif state.supplies[1] > office_data.mc_resupply_when_below then
				--status = "Working. Medical Supplies low. Assigned labourers will re-supply office."
				--supply_warning = supply_warning .. "Low on Medical Supplies."
				--table.insert(status_data,"Medical Supplies")
				if state.resupply == false then
					--status = "Working. Workcrew ordered to NOT re-supply office."
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
				
				--status = "Work halted. Medical Supplies needed."
				--supply_warning = "Medical Supplies required to heal wounded."
				table.insert(status_data,"Medical Supplies")
				if state.resupply == false then
					--status = "Work halted. Workcrew ordered to NOT re-supply office."
				else
					if state.buildingOwner then
						local ownername = query(state.buildingOwner,"getName")[1]
						local alertstring = "The Barbershop operated by " .. ownername .. " is out of Medical Supplies! Produce more to enable healing."
						
						send("rendCommandManager",
							"odinRendererStubMessage", --"odinRendererStubMessage",
							"ui\\commodityIcons.xml", -- iconskin
							"medical_bag", -- icon
							"Barbershop Needs Supplies", -- header text
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
				
				if state.resupply == false then
					SELF.tags.needs_resupply1 = nil
					SELF.tags.needs_resupply1_badly = nil
				else
					SELF.tags.needs_resupply1 = true
					SELF.tags.needs_resupply1_badly = true
				end
				
				SELF.tags.no_supplies1 = true
			end
			
			-- Sulphur tonic is a bonus, not a requirement.
			if state.supplies[2] >= office_data.lc_resupply_when_below then
				--status = status .. " Sulphur Tonic stocked."
				SELF.tags.no_supplies2 = nil
				SELF.tags.needs_resupply2 = nil
				SELF.tags.needs_resupply2_badly = nil
				
			elseif state.supplies[2] > office_data.mc_resupply_when_below then

				--status = status .. " Sulphur Tonic low."
				--supply_warning = supply_warning .. " Sulphur Tonic low."
				
				if state.resupply == false then
					SELF.tags.needs_resupply2 = nil
					SELF.tags.needs_resupply2_badly = nil
				else
					SELF.tags.needs_resupply2 = true
					SELF.tags.needs_resupply2_badly = nil
				end
				
				SELF.tags.no_supplies2 = nil
				
			elseif state.supplies[2] == 0 then

				--status = status .. " No Sulphur Tonic."
				table.insert(status_data,"Sulphur Tonic")
				--[[if state.resupply == false then
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
				end]]
				
				if state.resupply == false then
					SELF.tags.needs_resupply2 = nil
					SELF.tags.needs_resupply2_badly = nil
				else
					SELF.tags.needs_resupply2 = true
					SELF.tags.needs_resupply2_badly = true
				end
				
				SELF.tags.no_supplies2 = true
			end
			
			send("rendUIManager", "SetOfficeString", SELF, "noSuppliesWarning", supply_warning)
			send("rendUIManager", "SetOfficeString", SELF, "workPointsStatus", combined_warning_status(status_data))
		end
	>>

	state
	<<

	>>

	receive Create( stringstringMapHandle init )
	<<

	>>
	
	receive odinBuildingCompleteMessage( int handle, gameSimJobInstanceHandle ji )
	<<
		send("rendUIManager","SetOfficeString",SELF,"noChairWarning","Chairs needed!")
		send("rendUIManager", "SetOfficeInt", SELF, "chairsPresent",0)
		send("rendUIManager", "SetOfficeInt", SELF, "workPoints1", 0)
		send("rendUIManager", "SetOfficeInt", SELF, "workPoints2", 0)
		barbershop_reset_supply_text()
	>>
	
	receive InteractiveMessage( string messagereceived )
	<<
		printl("buildings", "barbershop received message: " .. messagereceived)
		if not state.completed or SELF.tags.slated_for_demolition then
			return
		end
		
		if messagereceived == "supply_on" then
			state.resupply = true
			barbershop_reset_supply_text()
		elseif messagereceived == "supply_off" then
			state.resupply = false
			barbershop_reset_supply_text()	
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
		barbershop_reset_supply_text()
	>>
	
	receive consumeSupplies( int tier, int count )
	<<
		state.supplies[tier] = state.supplies[tier] - count 
		send("rendUIManager", "SetOfficeInt", SELF, "workPoints" .. tostring(tier), state.supplies[tier])
		barbershop_reset_supply_text()
	>>
	
	receive recalculateQuality()
	<<
	
		-- count chairs
		local chair_count = 0
          local modules = query("gameWorkshopManager", "getBuildingModulesGameSide", SELF)[1]
		for k,v in pairs(modules) do 
               local module_tags = query(v, "getTags")[1]
               if module_tags.sittable then
				chair_count = chair_count + 1
			end
		end
		
		if chair_count > 0 then
			send("rendUIManager","SetOfficeString",SELF,"noChairWarning","")
		else
			send("rendUIManager","SetOfficeString",SELF,"noChairWarning","Chairs needed!")
		end
		send("rendUIManager", "SetOfficeInt", SELF, "chairsPresent",chair_count)	
     >>
	
	receive setBuildingOwner( gameObjectHandle newOwner )
	<<
		if newOwner and state.supplies[1] == 0 then
			
			local eventQ = query("gameSimEventManager",
								"startEvent",
								"supplies_warning_barbershop",
								{},
								{} )[1]
						
			send(eventQ,"registerBuilding",SELF)
		end
	>>
>>