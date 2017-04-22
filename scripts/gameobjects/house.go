gameobject "house" inherit "buildings"
<<
	local 
	<<
		function houseRefreshConditions()
			
			if SELF.tags.lower_class_house then
				local old_lc_cap_increase = state.lc_pop_cap_increase
				send("gameSession", "incSessionInt", "LcPopulationAllowed", old_lc_cap_increase * -1)
				send("gameSession", "incSessionInt", "totalPopulationAllowed", old_lc_cap_increase * -1)
				state.lc_pop_cap_increase = EntityDB.WorldStats.lc_house_pop_cap_increase
				
			elseif SELF.tags.middle_class_house then
				local old_mc_cap_increase = state.mc_pop_cap_increase
				send("gameSession", "incSessionInt", "McPopulationAllowed", old_mc_cap_increase * -1)
				send("gameSession", "incSessionInt", "totalPopulationAllowed", old_mc_cap_increase * -1)
				state.mc_pop_cap_increase = EntityDB.WorldStats.mc_house_pop_cap_increase
				
			end
				
			local modules = query("gameWorkshopManager", "getBuildingModulesGameSide", SELF)[1]
			
			-- Condition 1: 2x windows : adds +1
			-- Condition 2: 3x cots : adds +1
			-- Condition 3: 'great' or better quality: adds +1
			-- Condition 4: painting: adds +2
			
			local condition1 = false
			local condition2 = false
			local condition3 = false
			local condition4 = false
			
			if SELF.tags.lower_class_house then
				local window_count = 0
				local cot_count = 0
	
				for k,v in pairs(modules) do
					local tags = query(v,"getTags")[1]
					if not tags.module_disabled then
						if tags["window"] then
							window_count = window_count + 1
						elseif tags.lower_class_bed then
							cot_count = cot_count + 1
						elseif tags.painting then
							condition4 = true
						end
					end
				end
				
				if window_count >= 2 then
					condition1 = true
				end
				if cot_count >= 3 then
					condition2 = true
				end
				if state.buildingQuality >= 4 then
					condition3 = true
				end
				
			elseif SELF.tags.middle_class_house then
				local bed_count = 0
				local window_count = 0
				local tablechair_count = 0
				
				for k,v in pairs(modules) do
					local tags = query(v,"getTags")[1]
					if tags["window"] then
						window_count = window_count + 1
					elseif tags.middle_class_bed then
						bed_count = bed_count + 1
					elseif tags.middle_class_module and tags.chair and tags.table then
						tablechair_count = tablechair_count + 1
					elseif tags.painting then
						condition4 = true
					end
				end
				
				if bed_count >= 2 then
					condition1 = true
				end
				if window_count >= 2 and tablechair_count >= 2 then
					condition2 = true
				end
				if state.buildingQuality >= 6 then
					condition3 = true
				end
			end
				
			printl("buildings", "HouseRefreshConditions results: " ..
					tostring(condition1) .. " / " ..
					tostring(condition2) .. " / " ..
					tostring(condition3) .. " / " ..
					tostring(condition4) .. ") ")
				
			if SELF.tags.lower_class_house then
				if condition1 == true then
					state.lc_pop_cap_increase = state.lc_pop_cap_increase + 1
					send("rendUIManager", "SetOfficeString", SELF, "condition1todoText", "")-- "2x Windows")
					send("rendUIManager", "SetOfficeString", SELF, "condition1todoInt", "") -- "+1")
					send("rendUIManager", "SetOfficeString", SELF, "condition1DoneText", "2x Windows")
					send("rendUIManager", "SetOfficeString", SELF, "condition1DoneInt", "DONE +1")
				else
					send("rendUIManager", "SetOfficeString", SELF, "condition1todoText", "2x Windows")
					send("rendUIManager", "SetOfficeString", SELF, "condition1todoInt", "+1")
					send("rendUIManager", "SetOfficeString", SELF, "condition1DoneText", "")
					send("rendUIManager", "SetOfficeString", SELF, "condition1DoneInt", "")
				end
				
				if condition2 == true then
					state.lc_pop_cap_increase = state.lc_pop_cap_increase + 1
					send("rendUIManager", "SetOfficeString", SELF, "condition2todoText", "")
					send("rendUIManager", "SetOfficeString", SELF, "condition2todoInt", "") 
					send("rendUIManager", "SetOfficeString", SELF, "condition2DoneText", "3x Cots")
					send("rendUIManager", "SetOfficeString", SELF, "condition2DoneInt", "DONE +1")
				else
					send("rendUIManager", "SetOfficeString", SELF, "condition2todoText", "3x Cots")
					send("rendUIManager", "SetOfficeString", SELF, "condition2todoInt", "+1")
					send("rendUIManager", "SetOfficeString", SELF, "condition2DoneText", "")
					send("rendUIManager", "SetOfficeString", SELF, "condition2DoneInt", "")
				end
	
				if condition3 == true then
					state.lc_pop_cap_increase = state.lc_pop_cap_increase + 1
					send("rendUIManager", "SetOfficeString", SELF, "condition3todoText", "")
					send("rendUIManager", "SetOfficeString", SELF, "condition3todoInt", "") 
					send("rendUIManager", "SetOfficeString", SELF, "condition3DoneText", "Quality of +4 or better")
					send("rendUIManager", "SetOfficeString", SELF, "condition3DoneInt", "DONE +1")
				else
					send("rendUIManager", "SetOfficeString", SELF, "condition3todoText", "Quality of +4 or better")
					send("rendUIManager", "SetOfficeString", SELF, "condition3todoInt", "+1")
					send("rendUIManager", "SetOfficeString", SELF, "condition3DoneText", "")
					send("rendUIManager", "SetOfficeString", SELF, "condition3DoneInt", "")
				end

				if condition4 == true then
					state.lc_pop_cap_increase = state.lc_pop_cap_increase + 2
					send("rendUIManager", "SetOfficeString", SELF, "condition4todoText", "")
					send("rendUIManager", "SetOfficeString", SELF, "condition4todoInt", "") 
					send("rendUIManager", "SetOfficeString", SELF, "condition4DoneText", "1x Painting")
					send("rendUIManager", "SetOfficeString", SELF, "condition4DoneInt", "DONE +2")
				else
					send("rendUIManager", "SetOfficeString", SELF, "condition4todoText", "1x Painting")
					send("rendUIManager", "SetOfficeString", SELF, "condition4todoInt", "+2")
					send("rendUIManager", "SetOfficeString", SELF, "condition4DoneText", "")
					send("rendUIManager", "SetOfficeString", SELF, "condition4DoneInt", "")
				end
				
				send("rendUIManager", "SetOfficeString", SELF, "popCapBonusText1", tostring(state.lc_pop_cap_increase) )
				send("rendUIManager", "SetOfficeString", SELF, "popCapBonusText2", "7")
				
				send("gameSession", "incSessionInt", "LcPopulationAllowed", state.lc_pop_cap_increase )
				send("gameSession", "incSessionInt", "totalPopulationAllowed", state.lc_pop_cap_increase )
			
			elseif SELF.tags.middle_class_house then
				
				if condition1 == true then
					state.mc_pop_cap_increase = state.mc_pop_cap_increase + 1
					send("rendUIManager", "SetOfficeString", SELF, "condition1todoText", "")-- "2x Windows")
					send("rendUIManager", "SetOfficeString", SELF, "condition1todoInt", "") -- "+1")
					send("rendUIManager", "SetOfficeString", SELF, "condition1DoneText", "2x Practical Beds")
					send("rendUIManager", "SetOfficeString", SELF, "condition1DoneInt", "DONE +1")
				else
					send("rendUIManager", "SetOfficeString", SELF, "condition1todoText", "2x Practical Beds")
					send("rendUIManager", "SetOfficeString", SELF, "condition1todoInt", "+1")
					send("rendUIManager", "SetOfficeString", SELF, "condition1DoneText", "")
					send("rendUIManager", "SetOfficeString", SELF, "condition1DoneInt", "")
				end
				
				if condition2 == true then
					state.mc_pop_cap_increase = state.mc_pop_cap_increase + 1
					send("rendUIManager", "SetOfficeString", SELF, "condition2todoText", "")
					send("rendUIManager", "SetOfficeString", SELF, "condition2todoInt", "") 
					send("rendUIManager", "SetOfficeString", SELF, "condition2DoneText", "2x Windows + 2x Practical Table and Chair Set ")
					send("rendUIManager", "SetOfficeString", SELF, "condition2DoneInt", "DONE +1")
				else
					send("rendUIManager", "SetOfficeString", SELF, "condition2todoText", "2x Windows + 2x Practical Table and Chair Set ")
					send("rendUIManager", "SetOfficeString", SELF, "condition2todoInt", "+1")
					send("rendUIManager", "SetOfficeString", SELF, "condition2DoneText", "")
					send("rendUIManager", "SetOfficeString", SELF, "condition2DoneInt", "")
				end
	
				if condition3 == true then
					state.mc_pop_cap_increase = state.mc_pop_cap_increase + 1
					send("rendUIManager", "SetOfficeString", SELF, "condition3todoText", "")
					send("rendUIManager", "SetOfficeString", SELF, "condition3todoInt", "") 
					send("rendUIManager", "SetOfficeString", SELF, "condition3DoneText", "Quality of +6 or better")
					send("rendUIManager", "SetOfficeString", SELF, "condition3DoneInt", "DONE +1")
				else
					send("rendUIManager", "SetOfficeString", SELF, "condition3todoText", "Quality of +6 or better")
					send("rendUIManager", "SetOfficeString", SELF, "condition3todoInt", "+1")
					send("rendUIManager", "SetOfficeString", SELF, "condition3DoneText", "")
					send("rendUIManager", "SetOfficeString", SELF, "condition3DoneInt", "")
				end

				if condition4 == true then
					state.mc_pop_cap_increase = state.mc_pop_cap_increase + 1
					send("rendUIManager", "SetOfficeString", SELF, "condition4todoText", "")
					send("rendUIManager", "SetOfficeString", SELF, "condition4todoInt", "") 
					send("rendUIManager", "SetOfficeString", SELF, "condition4DoneText", "1x Painting")
					send("rendUIManager", "SetOfficeString", SELF, "condition4DoneInt", "DONE +1")
				else
					send("rendUIManager", "SetOfficeString", SELF, "condition4todoText", "1x Painting")
					send("rendUIManager", "SetOfficeString", SELF, "condition4todoInt", "+1") 
					send("rendUIManager", "SetOfficeString", SELF, "condition4DoneText", "")
					send("rendUIManager", "SetOfficeString", SELF, "condition4DoneInt", "")
				end
				
				send("rendUIManager", "SetOfficeString", SELF, "popCapBonusText1", tostring(state.mc_pop_cap_increase) )
				send("rendUIManager", "SetOfficeString", SELF, "popCapBonusText2", "5")
				
				send("gameSession", "incSessionInt", "McPopulationAllowed", state.mc_pop_cap_increase )
				send("gameSession", "incSessionInt", "totalPopulationAllowed", state.mc_pop_cap_increase )
			end
		end
	>>

	state
	<<
		int rOH
		string buildingName
		int buildingQuality
		string buildingQualityName
		table squares
		bool completed
		
		table materials
		table parent

		bool claimed
		gameObjectHandle buildingOwner
		gameSimAssignmentHandle curConstructionAssignment
	>>

	receive Create( stringstringMapHandle init )
	<<
		state.buildingTitle = "House"
		if init.legacyString then
			if EntityDB[init.legacyString].displayName then
				state.buildingTitle = EntityDB[init.legacyString].displayName
			end
		end
		
		--[[if state.displayName then
			send("rendUIManager", "SetOfficeString", SELF, "houseType", state.buildingTitle)
		else
			send("rendUIManager", "SetOfficeString", SELF, "houseType", state.displayName)
		end]]
		
		state.lc_pop_cap_increase = 0
		state.mc_pop_cap_increase = 0
		state.uc_pop_cap_increase = 0
	>>
	
	receive InteractiveMessage( string messagereceived )
	<<
		printl("buildings","house.go : " .. tostring(SELF.id) .. " / Message Received: " .. messagereceived )
		if not state.completed then
			return
		end

		-- OTHER STUFF HERE
		if messagereceived == "build" then
			-- Open menu
			send("rendUIManager", "OpenHousingMenu", SELF, state.buildingName)
			if state.buildingOwner == nil then
				send("rendUIManager", "SetHousingMenuHeader", SELF, state.buildingName, "carpentry_icon");
			else
				send("rendUIManager", "SetHousingMenuHeader", SELF, state.buildingName, "carpentry_icon");
			end
		end
	>>


	receive odinBuildingCompleteMessage ( int handle, gameSimJobInstanceHandle ji )
	<<
		send("gameWorkshopManager", "AddOffice", SELF, state.buildingName)
		
		if SELF.tags.lower_class_house then
			state.residents = {}
			SELF.tags.open_lc_house_slots = true
			
			send("rendUIManager",
				"SetOfficeString",
				SELF,
				"houseDescription",
				"A home for labourers. Adds base +2 Labourer population capacity." )

			--send("rendUIManager","SetOfficeString",SELF,"houseType","Low Class Bunkhouse" )
				
			send("rendUIManager",
				"SetOfficeString",
				SELF,
				"popCapBonusTitle",
				"Labourer Population Cap Bonus: " )
			
			send("rendUIManager",
				"SetOfficeString",
				SELF,
				"popcapTooltip",
				"Adds to the maximum number of Labourers." )
			
		elseif SELF.tags.middle_class_house then
			send("rendUIManager",
				"SetOfficeString",
				SELF,
				"houseDescription",
				"A home for overseers. Adds base +1 Overseer population capacity." )
			
			--send("rendUIManager","SetOfficeString",SELF,"houseType","Middle Class House" )
			send("rendUIManager",
				"SetOfficeString",
				SELF,
				"popCapBonusTitle",
				"Overseer Population Cap Bonus: " )
			
			send("rendUIManager",
				"SetOfficeString",
				SELF,
				"popcapTooltip",
				"Adds to the maximum number of Overseers." )
			
		elseif SELF.tags.upper_class_house then
			
			send("rendUIManager",
				"SetOfficeString",
				SELF,
				"houseDescription",
				"A home for aristocrats.")
			
				--"A home for aristocrats. Adds a base of +1 Upper Class population capacity to the colony." )
			--send("rendUIManager","SetOfficeString",SELF,"houseType","Upper Class Manor" )
			send("rendUIManager",
				"SetOfficeString",
				SELF,
				"popCapBonusTitle",
			     " ")
				--"Upper Class Population Cap Bonus: " )
				
				send("rendUIManager","SetOfficeString",SELF,"popcapTooltip","Aristocrats care not for immigration regulations." )
				
				-- display nothing for now.
				send("rendUIManager", "SetOfficeString", SELF, "popCapBonusText1"," " )
				send("rendUIManager", "SetOfficeString", SELF, "popCapBonusText2"," ")
				
			send("rendUIManager", "SetOfficeString", SELF, "condition1todoText", "")
			send("rendUIManager", "SetOfficeString", SELF, "condition1todoInt", "") 
			send("rendUIManager", "SetOfficeString", SELF, "condition1DoneText", "")
			send("rendUIManager", "SetOfficeString", SELF, "condition1DoneInt", "")
			send("rendUIManager", "SetOfficeString", SELF, "condition2todoText", "")
			send("rendUIManager", "SetOfficeString", SELF, "condition2todoInt", "") 
			send("rendUIManager", "SetOfficeString", SELF, "condition2DoneText", "")
			send("rendUIManager", "SetOfficeString", SELF, "condition2DoneInt", "")
			send("rendUIManager", "SetOfficeString", SELF, "condition3todoText", "")
			send("rendUIManager", "SetOfficeString", SELF, "condition3todoInt", "") 
			send("rendUIManager", "SetOfficeString", SELF, "condition3DoneText", "")
			send("rendUIManager", "SetOfficeString", SELF, "condition3DoneInt", "")
			send("rendUIManager", "SetOfficeString", SELF, "condition4todoText", "")
			send("rendUIManager", "SetOfficeString", SELF, "condition4todoInt", "") 
			send("rendUIManager", "SetOfficeString", SELF, "condition4DoneText", "")
			send("rendUIManager", "SetOfficeString", SELF, "condition4DoneInt", "")
		end
		
		
		send("gameWorkshopManager", "AddHouse", SELF, state.buildingName)

		if state.buildingName == "Upper Class House" then
			send("gameSession", "incSessionInt", "UCHousesProduced", 1)
		end
		
		if query("gameSession", "getSessionBool", "enableContextualTutorials")[1] == true then
			if not query("gameSession", "getSessionBool", "qualityTutorialDone")[1] then
				local eventQ = query("gameSimEventManager", "startEvent", "tutorial_quality", {}, {})
			end
		end

		ready()
		
		send("rendCommandManager",
			"odinRendererTickerMessage",
			"A " .. state.buildingName .. " finished construction on day " ..
				query("gameSession","getSessionInt","dayCount")[1] .. ".",
			"housing",
			"ui\\orderIcons.xml")
		
		state.hitpoints =  #state.squares --10 -- you get your hitpoints!
		state.hitpointsmax = #state.squares -- For display and possibly repair
		send("rendUIManager", "SetOfficeInt", SELF, "buildingHP", state.hitpoints)
		send("rendUIManager", "SetOfficeInt", SELF, "buildingHPMax", state.hitpointsmax)
		send("rendUIManager", "SetOfficeString", SELF, "buildingHPDescription", "Undamaged")
		send("rendUIManager","SetOfficeString",SELF,"buildingFancyName",state.buildingName )
		
		send(SELF, "recalculateQuality")
	>>

	receive Update()
	<<

	>>

	receive recalculateQuality()
	<<
		houseRefreshConditions()
	>>

	respond registerAsOccupant( gameObjectHandle character )
	<<
		
	>>
>>