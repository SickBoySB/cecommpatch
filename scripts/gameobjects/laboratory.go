gameobject "laboratory" inherit "office"
<<
	local 
	<<		
		function laboratory_reset_supply_text()
			
			local status = ""
			local supply_warning = ""
			local office_data = EntityDB[ state.entityName ]

			if state.supplies[1] >= office_data.lc_resupply_when_below then

				--status = "Working. Stocked with supplies."
				if state.resupply == false then
					--status = "Working. Resupply HALTED."
				end

				SELF.tags.no_supplies1 = nil
				SELF.tags.needs_resupply1 = nil
				SELF.tags.needs_resupply1_badly = nil
				
			elseif state.supplies[1] > office_data.mc_resupply_when_below then
	
				--status = "Working. Supplies low."
				status = "Science Materials (low) needed!"
				supply_warning = ""
				if state.resupply == false then
					--status = "Working. Resupply HALTED."
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
			
				--status = "Work HALTED. Supplies needed."
				supply_warning = ""
				status = "Science Materials needed!"
				if state.resupply == false then
					--status = "Work HALTED. Resupply HALTED."
				else
					if state.buildingOwner then
						local ownername = query(state.buildingOwner,"getName")[1]
						local alertstring = "The Laboratory operated by " .. ownername .. " is out of Science Materials! Produce more to enable research."
						
						send("rendCommandManager",
							"odinRendererStubMessage", --"odinRendererStubMessage",
							"ui\\orderIcons.xml", -- iconskin
							"laboratory", -- icon
							"Laboratory Needs Supplies", -- header text
							alertstring, -- text description
							"Left-click to zoom. Right-click to dismiss.", -- action string
							"laboratoryProblem", -- alert typasde (for stacking)
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
			
			if not state.buildingOwner then
				--status = "Work HALTED. Overseer needed."
			end
			
			send("rendUIManager", "SetOfficeString", SELF, "noSuppliesWarning",supply_warning)
			send("rendUIManager", "SetOfficeString", SELF, "workPointsStatus", status)
		end
		
		
		function setupNewProject ( program, slotNumber )
			if not state.availableProjects[program] then
				state.availableProjects[program] = {}
			end
			table.insert(state.availableProjects[program], slotNumber, {})
		end

		function updateOfficeVars( name )
			
			--printl("buildings", "laboratory doing updateOfficeVars for: " .. name )
	
			send("rendUIManager", "SetOfficeInt", SELF, "numAvailableProjects" .. name, #state.availableProjects[name])
			send("rendUIManager", "SetOfficeInt", SELF, "researchPoints" .. name, state.researchPoints[name])

			-- clear existing stuff
			--[[for i=1,state.maxAvailableProjects do
				send("rendUIManager", "SetOfficeString", SELF, "ProjectName" .. name .. i, "")
				send("rendUIManager", "SetOfficeString", SELF, "ProjectDisplayName" .. name .. i, "")
				send("rendUIManager", "SetOfficeString", SELF, "ProjectIcon" ..name .. name .. i, "")
				send("rendUIManager", "SetOfficeString", SELF, "ProjectIconSkin" .. name .. i, "")
				send("rendUIManager", "SetOfficeInt", SELF, "ProjectCost" .. name .. i, -1)
				send("rendUIManager", "SetOfficeBool", SELF, "ProjectEnabled" .. name .. i, false)
			end]]
			
			local count = 1
			for k2,v2 in pairs ( state.availableProjects[name] ) do
				send("rendUIManager", "SetOfficeString", SELF, name .. "ProjectName" .. count, v2.name)
				send("rendUIManager", "SetOfficeString", SELF, name .. "ProjectDisplayName" .. count, v2.displayName)
				send("rendUIManager", "SetOfficeString", SELF, name .. "ProjectDesc" .. count, v2.description)
				send("rendUIManager", "SetOfficeString", SELF, name .. "ProjectIcon" .. count, v2.icon)
				send("rendUIManager", "SetOfficeString", SELF, name .. "ProjectIconSkin" .. count, v2.iconSkin)
				send("rendUIManager", "SetOfficeInt", SELF, name .. "ProjectCost" .. count, v2.cost)
				send("rendUIManager", "SetOfficeBool", SELF, name .. "ProjectEnabled" .. count, v2.enabled)
				count = count + 1
			end
		end
	>>

	state
	<<
		int maxAvailableProjects
		int researchPoolCap
		int numUnlockedSlots
		table availableProjects
		table researchPoints
		string currentProjectGroup
		table numAvailableProjects
		table allPrograms
		table completedProjects
		bool chalkboardPresent
		bool macroscopePresent
		bool vacuumChamberPresent
	>>

	receive Create( stringstringMapHandle init )
	<<
		state.researchPoolCap = 20
		state.maxAvailableProjects = 3
		state.numUnlockedSlots = 3 
		state.allPrograms = {}
          state.scienceGain = 0
		
		state.chalkboardPresent = false
		state.macroscopePresent = false
		state.vacuumChamberPresent = false	
	>>

	receive odinBuildingCompleteMessage( int handle, gameSimJobInstanceHandle ji )
	<<
		-- built table of all valid program names
		for name,data in pairs( EntityDB.researchTable.researchPrograms ) do
			state.allPrograms[#state.allPrograms+1] = name
		end
		
		-- default area of focus.
		state.currentProjectGroup = "colonyProduction"

		-- number of projects available per program
		state.numAvailableProjects = {}
		state.availableProjects = {}

		for k, programName in pairs(state.allPrograms) do
				
			printl("buildings", "laboratory " .. tostring(SELF.id) .. " initiatizing program: " .. programName )
			
			state.numAvailableProjects[programName] = 1
			state.availableProjects[programName] = {}
			state.researchPoints[programName] = 0
			
			setupNewProject( programName, 1)
			
			send("rendUIManager",
				"SetOfficeInt",
				SELF,
				"researchPoints" .. programName,
				state.researchPoints[programName])
			
			state.completedProjects[programName] = {}
		end
		
		for k, programName  in pairs(state.allPrograms) do
			send(SELF, "generateProjects", programName, nil)
		end
		
		send("rendUIManager","SetOfficeString",SELF,"scienceSkillDescriptor", "n/a")
		
		send(SELF,"refreshSkillDisplay")
          send(SELF,"labRefreshConditions")
		
		for k,v in pairs(state.allPrograms) do
			updateOfficeVars( v )
		end
		
		
		send("rendUIManager","SetOfficeInt",SELF,"workPoints1",0)
		--send("rendUIManager","SetOfficeInt",SELF,"workPoints2",0)
		--send("rendUIManager","SetOfficeInt",SELF,"workPoints3",0)
		
		send("rendUIManager", "SetOfficeString", SELF, "workPointsStatus", "Science Materials needed!")
		send("rendUIManager","SetOfficeString",SELF,"noEquipmentWarning","Research module needed!")
		
		laboratory_reset_supply_text()
	>>

	receive InteractiveMessage( string messagereceived )
	<<
		if not state.completed then
			return
		end
		
		if messagereceived == "supply_on" then
			state.resupply = true
			laboratory_reset_supply_text()
			
		elseif messagereceived == "supply_off" then
			state.resupply = false
			laboratory_reset_supply_text()
			
		else
			-- if string has "Reset" in it, extract the reset and see if the remainder matches a programName
			local resetFindResult = string.find(messagereceived, "Reset")
			if resetFindResult then
				local programName = string.sub( messagereceived, 1, resetFindResult-1 )
				
				for k,v in pairs( state.allPrograms ) do
					if v == programName then 
						-- we're in business!
						if state.researchPoints[ programName ] > 2 then
							state.researchPoints[ programName ] = state.researchPoints[ programName ] - 3
							send("rendUIManager", "SetOfficeInt", SELF, "researchPoints" .. programName, state.researchPoints[programName])
							send(SELF, "generateProjects", programName, nil)
	
							updateOfficeVars( programName )
							break
						end
					end
				end
				
				return
			end
			
			-- if string has "ResearchProgram_" in it, do the stuff for that program.
			local researchProgramFindResult = string.find(messagereceived, "ResearchProgram_")
			if researchProgramFindResult then
				
				-- TODO: safetly validation here?
				
				local programName = string.sub( messagereceived, 17 )
				printl("buildings", "laboratory (ResearchProgram_) found substring in messagereceived: " .. programName )
			
				local oldProgramName = state.currentProjectGroup -- pass this in so we only update the old & new groups
				state.currentProjectGroup = programName
				
				-- set enabled to false in prev. program
				-- TODO use: number state.maxAvailableProjects
				
				send("rendUIManager", "SetOfficeInt", SELF,
					"numAvailableProjects" .. oldProgramName, #state.availableProjects[oldProgramName])
				
				send("rendUIManager", "SetOfficeInt", SELF,
					"researchPoints" .. oldProgramName, state.researchPoints[oldProgramName])
				
				-- CECOMMPATCH - Not sure why this was being triggered... removing fixes "project completed" bug
				--for i=1,3 do
					--state.availableProjects[oldProgramName][i].enabled = false
				--end
				
				local count = 1
				for k,v in pairs( state.availableProjects[oldProgramName] ) do
					--v.enabled = false

					send("rendUIManager", "SetOfficeString", SELF, oldProgramName .. "ProjectName" .. count, v.name)
					send("rendUIManager", "SetOfficeString", SELF, oldProgramName .. "ProjectDisplayName" .. count, v.displayName)
					send("rendUIManager", "SetOfficeString", SELF, oldProgramName .. "ProjectDesc" .. count, v.description)
					send("rendUIManager", "SetOfficeString", SELF, oldProgramName .. "ProjectIcon" .. count, v.icon)
					send("rendUIManager", "SetOfficeString", SELF, oldProgramName .. "ProjectIconSkin" .. count, v.iconSkin)
					send("rendUIManager", "SetOfficeInt",    SELF, oldProgramName .. "ProjectCost" .. count, v.cost)
					send("rendUIManager", "SetOfficeBool",   SELF, oldProgramName .. "ProjectEnabled" .. count, v.enabled)
					count = count + 1
				end
				
				
				-- set enabled to true in new program where points are high enough to enable proj.
				
				send("rendUIManager", "SetOfficeInt", SELF,
					"numAvailableProjects" .. programName, #state.availableProjects[programName])
				
				send("rendUIManager", "SetOfficeInt", SELF,
					"researchPoints" .. programName, state.researchPoints[programName])
				
				for i=1,3 do
					state.availableProjects[programName][i].enabled = true
				end
				
				local count = 1
				for k,v in pairs ( state.availableProjects[programName] ) do
					--[[if state.researchPoints[programName] >= v.cost then
						v.enabled = true
					end]]
					
					send("rendUIManager", "SetOfficeString", SELF, programName .. "ProjectName" .. count, v.name)
					send("rendUIManager", "SetOfficeString", SELF, programName .. "ProjectDisplayName" .. count, v.displayName)
					send("rendUIManager", "SetOfficeString", SELF, programName .. "ProjectDesc" .. count, v.description)
					send("rendUIManager", "SetOfficeString", SELF, programName .. "ProjectIcon" .. count, v.icon)
					send("rendUIManager", "SetOfficeString", SELF, programName .. "ProjectIconSkin" .. count, v.iconSkin)
					send("rendUIManager", "SetOfficeInt", 	 SELF, programName .. "ProjectCost" .. count, v.cost)
					send("rendUIManager", "SetOfficeBool",   SELF, programName .. "ProjectEnabled" .. count, v.enabled)
					count = count + 1
				end
				
				printl("buildings", "laboratory: UI: updating program: " .. oldProgramName) 
				
				return
			end
			
			-- if string has "project_" in it, do the stuff for that project.
			local projectFindResult = string.find(messagereceived, "project_")
			if projectFindResult then
				
				local projectName = string.sub( messagereceived, projectFindResult )
				
				local programName = string.sub( messagereceived, 9 )
				local slotNumber = tonumber( string.sub( programName, -1) )
				local programName = string.sub( programName, 1, -2 )-- for real.
				
				-- TODO: add case/check for if another lab already discovered the project while this was being studied?
				
				local selectedProject = state.availableProjects[programName][ slotNumber ]
				local sessionVarName = selectedProject.sessionVarName
				
				if selectedProject.eventName then
					local eventQ = query("gameSimEventManager",
									 "startEvent",
									 selectedProject.eventName,
									 {},
									 {})
				end
				
				state.researchPoints[programName] = state.researchPoints[programName] - selectedProject.cost
				
				send("rendUIManager",
					"SetOfficeInt",
					SELF,
					"researchPoints" .. programName,
					state.researchPoints[programName])
				
				table.insert( state.completedProjects[programName],selectedProject.name )
				
				if sessionVarName then
					send("gameSession", "setSessionBool", sessionVarName, true)
				end
				
				send(SELF,"generateProjects",programName,slotNumber)
				-- only one program group.
				updateOfficeVars( programName )
				
				return
			end
		end	
	>>

	receive IncrementResearchPoints( gameObjectHandle researcher )
	<<
		if state.researchPoints[ state.currentProjectGroup ] and
			 state.researchPoints[ state.currentProjectGroup ] >= state.researchPoolCap then
			
			if researcher then
				local btags = query(researcher, "getTags")[1]
				if btags.middle_class then
					
					local name = query(researcher,"getName")[1]
					local pronoun = query(researcher,"getPronoun")[3]
					
					local tickerText = name .. " 's laboratory is completely filled with Science and can't fit any more! \z
						Spend some research points, or " .. name .. " will just keep working on " .. pronoun .. " curriculum vitae \z
						instead of doing Science."
					
					send("rendCommandManager",
						"odinRendererFYIMessage",
						"ui\\orderIcons.xml",
						"laboratory",
						"Science Points Full!", -- header text
						tickerText, -- text description
						"Left-click for details. Right-click to dismiss.", -- action string
						"laboratoryFull", -- alert type (for stacking)
						"ui//eventart//doing_science.png", -- imagename for bg
						"high", -- importance: low / high / critical
						state.rOH, -- object ID
						30 * 1000, -- duration in ms
						0, -- "snooze" time if triggered multiple times in rapid succession
						nil) -- gameobjecthandle of director, null if none
				end
			end
			
		else
			state.researchPoints[state.currentProjectGroup] = state.researchPoints[state.currentProjectGroup] + state.scienceGain
		end
		
		send("rendUIManager", "SetOfficeInt", SELF, "researchPoints" .. state.currentProjectGroup, state.researchPoints[state.currentProjectGroup])
		
		if state.currentProjectGroup == "agriculture" then
			send("rendUIManager","SetOfficeString",SELF,"agriculturePoints", state.researchPoints[state.currentProjectGroup])
		end
	>>
	
	receive generateProjects(string program, int slotnumber )
	<<
		function findValidProject()
			
			local projectSet = EntityDB.researchTable.researchPrograms[ program ]
			local done = false
			
			while not done do
				-- find and test completely random entry until valid entry is found.
				local r = rand(1,#projectSet)
				local project = projectSet[r]
				
				-- test project
				
			end
			
			return 
		end
	
		if not slotnumber then
			printl("buildings", "laboratory doing generateProjects for: " .. program )
			for i=1,state.maxAvailableProjects do
				if state.availableProjects[program][i] then
					state.availableProjects[program][i] = nil
				end 
			end
		else
			printl("buildings", "laboratory doing generateNewProjects for: " .. program .. ", slotnumber = " .. slotnumber )
			state.availableProjects[program][slotnumber] = nil
		end
		
		
		local projectSet = EntityDB.researchTable.researchPrograms[ program ]
		local possibleProjects = {}
		
		for k,project in pairs( projectSet ) do
			-- project researched?
			if not query("gameSession","getSessionBool",project.sessionVarName)[1] then
				-- check unlock requirement.
				if not project.requireSessionBool or
					(project.requireSessionBool and query("gameSession","getSessionBool", project.requireSessionBool)[1]) then
					
					-- check if biome matches, if biome-locked.
					if not project.biomes or
						(project.biomes and	project.biomes[ query("gameSession","getSessionString","biome")[1] ] ) then
						
						if slotnumber then
							-- check vs. existing slots and make sure this isn't one of them.
							local valid = true
							for slot,object in pairs(state.availableProjects[program]) do
								if object.name == project.name then
									-- no
									valid = false
									break
								end
							end
							
							if valid then
								possibleProjects[#possibleProjects + 1] = project
							end
						else
							possibleProjects[#possibleProjects + 1] = project
						end
					end
				end
			end
		end
		
		-- if options exists, put random project in THEN remove from list of possible projects
		if slotnumber then
			if #possibleProjects == 0 then
				state.availableProjects[program][slotnumber] = {
					cost = 999,
					name = "no_research",
					displayName = "No More Research Available",
					description = "You researched it all, your scientists are empty husks, bereft of ideas.",
					iconSkin = "ui//orderIcons.xml",
					icon = "x_black_icon",
					enabled = false, }
				
			else
				state.availableProjects[program][slotnumber] = possibleProjects[ rand(1,#possibleProjects) ]
				state.availableProjects[program][slotnumber].enabled = true
			end
		else 
			for i=1,state.maxAvailableProjects do
				if #possibleProjects == 0 then
					state.availableProjects[program][i] = {
						cost = 999,
						name = "no_research",
						displayName = "No More Research Available",
						description = "You researched it all, your scientists are empty husks, bereft of ideas.",
						iconSkin = "ui//orderIcons.xml",
						icon = "x_black_icon",
						enabled = false, }
					
				else
					-- add it! (and remove chosen project from list of projects.)
					state.availableProjects[program][i] = table.remove( possibleProjects, rand(1,#possibleProjects) )
					state.availableProjects[program][i].enabled = true
				end
			end
		end
	>>
	
	receive setBuildingOwner(gameObjectHandle newOwner)
	<<
		send(SELF,"refreshSkillDisplay")
		laboratory_reset_supply_text()

		if newOwner and state.supplies[1] == 0 then
			
			local eventQ = query("gameSimEventManager",
								"startEvent",
								"supplies_warning_laboratory",
								{},
								{} )[1]
						
			send(eventQ,"registerBuilding",SELF)
		end
	>>
	
	receive refreshSkillDisplay()
	<<
		-- owner skilled up OR was added
		if state.buildingOwner then 
			local skill = query(state.buildingOwner,"getEffectiveSkillLevel","science")[1]
			
			state.researchPoolCap = skill * 20
			
			printl("buildings", "Laboratory " .. tostring(SELF.id) .. " got owner with science skill = " .. skill .. " resulting in science points cap of: " .. state.researchPoolCap )
			send("rendUIManager","SetOfficeString",SELF,"scienceSkillDescriptor", EntityDB.HumanStats.skillLevelStrings[skill + 1] )
		else
			state.researchPoolCap = 20
			printl("buildings", "Laboratory " .. tostring(SELF.id) .. " has NO owner result, state.researchPoolCap in science points cap of: " .. state.researchPoolCap )
			send("rendUIManager","SetOfficeString",SELF,"scienceSkillDescriptor", "n/a")
		end
		
		if state.researchPoints.agriculture > state.researchPoolCap then state.researchPoints.agriculture = state.researchPoolCap end
		if state.researchPoints.colonyProduction > state.researchPoolCap then state.researchPoints.colonyProduction = state.researchPoolCap end
		if state.researchPoints.military > state.researchPoolCap then state.researchPoints.military = state.researchPoolCap end
		
		for k,v in pairs(state.allPrograms) do
			updateOfficeVars( v )
		end
		
		send("rendUIManager","SetOfficeInt",SELF,"agriculturePointsMax", state.researchPoolCap)
		send("rendUIManager","SetOfficeInt",SELF,"colonyProductionPointsMax", state.researchPoolCap)
		send("rendUIManager","SetOfficeInt",SELF,"militaryPointsMax", state.researchPoolCap)
		
		--send("rendUIManager","SetOfficeInt",SELF,"sciencePointsMax", state.researchPoolCap) --
		--send("rendUIManager","SetOfficeInt",SELF,"miningMetallurgyPointsNumberContainer", tostring(state.researchPoolCap))
		--send("rendUIManager","SetOfficeInt",SELF,"militaryPointsNumberContainer", tostring(state.researchPoolCap))
		--send("rendUIManager","SetOfficeInt",SELF,"miningMetallurgyPointsNumberContainer", tostring(state.researchPoolCap))
     >>
     
     receive labRefreshConditions()
     <<
          local scienceGainPoints = 0
          local modules = query("gameWorkshopManager", "getBuildingModulesGameSide", SELF)[1]
          
          -- Condition 1: 2x windows : adds +1
          -- Condition 2: 3x cots : adds +1
          -- Condition 3: 'great' or better quality: adds +1
          -- Condition 4: painting: adds +2
          local condition1 = false --Chalkboard
          local condition2 = false --Macroscope
          local condition3 = false --Vacuum Chamber
          local condition4 = false --1x Barometer, 2x MC Window

          local window_count = 0
          local barometer_count = 0

          for k,v in pairs(modules) do
               local tags = query(v,"getTags")[1]
               if not tags.module_disabled then
                    if tags["window"] then
                         window_count = window_count + 1
                    elseif tags.chalkboard then
                         condition1 = true
                    elseif tags.vacuum_chamber then
                         condition3 = true
                    elseif tags.macroscope then
                         condition2 = true
                    elseif tags.barometer then
                         barometer_count = barometer_count + 1
                    end
				--Note to self: Maybe something for plaques here later?
               end
          end
          
          if (window_count >= 2) and (barometer_count >= 1) then
               condition4 = true
          end

          if condition1 == true then
               scienceGainPoints = scienceGainPoints + 1
               send("rendUIManager", "SetOfficeString", SELF, "condition1todoText", "")-- "2x Windows")
               send("rendUIManager", "SetOfficeString", SELF, "condition1todoInt", "") -- "+1")
               send("rendUIManager", "SetOfficeString", SELF, "condition1DoneText", "Build Chalkboard")
               send("rendUIManager", "SetOfficeString", SELF, "condition1DoneInt", "DONE +1")
          else
               send("rendUIManager", "SetOfficeString", SELF, "condition1todoText", "Build Chalkboard")
               send("rendUIManager", "SetOfficeString", SELF, "condition1todoInt", "+1")
               send("rendUIManager", "SetOfficeString", SELF, "condition1DoneText", "")
               send("rendUIManager", "SetOfficeString", SELF, "condition1DoneInt", "")
          end
          
          if condition2 == true then
               scienceGainPoints = scienceGainPoints + 1
               send("rendUIManager", "SetOfficeString", SELF, "condition2todoText", "")
               send("rendUIManager", "SetOfficeString", SELF, "condition2todoInt", "") 
               send("rendUIManager", "SetOfficeString", SELF, "condition2DoneText", "Build Macroscope")
               send("rendUIManager", "SetOfficeString", SELF, "condition2DoneInt", "DONE +1")
          else
               send("rendUIManager", "SetOfficeString", SELF, "condition2todoText", "Build Macroscope")
               send("rendUIManager", "SetOfficeString", SELF, "condition2todoInt", "+1")
               send("rendUIManager", "SetOfficeString", SELF, "condition2DoneText", "")
               send("rendUIManager", "SetOfficeString", SELF, "condition2DoneInt", "")
          end

          if condition3 == true then
               scienceGainPoints = scienceGainPoints + 2
               send("rendUIManager", "SetOfficeString", SELF, "condition3todoText", "")
               send("rendUIManager", "SetOfficeString", SELF, "condition3todoInt", "") 
               send("rendUIManager", "SetOfficeString", SELF, "condition3DoneText", "Build Vacuum Chamber")
               send("rendUIManager", "SetOfficeString", SELF, "condition3DoneInt", "DONE +2")
          else
               send("rendUIManager", "SetOfficeString", SELF, "condition3todoText", "Build Vacuum Chamber")
               send("rendUIManager", "SetOfficeString", SELF, "condition3todoInt", "+2")
               send("rendUIManager", "SetOfficeString", SELF, "condition3DoneText", "")
               send("rendUIManager", "SetOfficeString", SELF, "condition3DoneInt", "")
          end
               
          if condition4 == true then
               scienceGainPoints = scienceGainPoints + 1
               send("rendUIManager", "SetOfficeString", SELF, "condition4todoText", "")
               send("rendUIManager", "SetOfficeString", SELF, "condition4todoInt", "") 
               send("rendUIManager", "SetOfficeString", SELF, "condition4DoneText", "2x Window, 1x Barometer")
               send("rendUIManager", "SetOfficeString", SELF, "condition4DoneInt", "DONE +1")
          else
               send("rendUIManager", "SetOfficeString", SELF, "condition4todoText", "2x Window, 1x Barometer")
               send("rendUIManager", "SetOfficeString", SELF, "condition4todoInt", "+1")
               send("rendUIManager", "SetOfficeString", SELF, "condition4DoneText", "")
               send("rendUIManager", "SetOfficeString", SELF, "condition4DoneInt", "")
          end
          
          state.scienceGain = scienceGainPoints
          send("rendUIManager", "SetOfficeString", SELF, "scienceGainPointsText1", tostring(state.scienceGain) )
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
		laboratory_reset_supply_text()
	>>
	
	receive consumeSupplies( int tier, int count )
	<<
		state.supplies[tier] = state.supplies[tier] - count 
		send("rendUIManager", "SetOfficeInt", SELF, "workPoints" .. tostring(tier), state.supplies[tier])
		laboratory_reset_supply_text()
	>>
	
	receive recalculateQuality()
	<<
		-- count lab_equipment
		local equipment_count = 0
          local modules = query("gameWorkshopManager", "getBuildingModulesGameSide", SELF)[1]
		for k,v in pairs(modules) do 
               local module_tags = query(v, "getTags")[1]
               if module_tags.lab_equipment then
				equipment_count = equipment_count + 1
			end
		end
		
		if equipment_count > 0 then
			send("rendUIManager","SetOfficeString",SELF,"noEquipmentWarning","")
		else
			send("rendUIManager","SetOfficeString",SELF,"noEquipmentWarning","Research module needed!")
		end
		
		send(SELF, "labRefreshConditions")
     >>
>>