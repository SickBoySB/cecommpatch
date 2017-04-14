gameobject "buildings"
<<
	local
	<<
	
		function combined_warning_status( warning_data )
			-- CECOMMPATCH function. Used to combine warnings/status text for offices
			local o = ""
			
			if #warning_data > 0 then
				if #warning_data == 2 then
					o = table.concat(warning_data," and ")
				else
					o = table.concat(warning_data,", ")
				end
				
				o = o .. " needed!"
			end
			
			return o
		end

	>>
	
	state
	<<
		bool slatedForDemolition
		gameSimAssignmentHandle curConstructionAssignment
		gameSimAssignmentHandle curDestructionAssignment
		table modules
		table points
		table squares
		table contents
		
		int buildingUpgradeLevel
		int buildingQuality
		string buildingQualityName
		int hitpoints
		int hitpointsmax
		
		bool repair_material_dropped_off
		table repair_materials
	>>

	receive Create( stringstringMapHandle init )
	<<
		printl("buildings", "received create for " .. tostring(init.legacyString) .. ", making " .. tostring(SELF.id) )
	
		local myName = init["legacyString"]
		state.buildingName = myName
		state.buildingQuality = 0
		state.buildingQualityName = "Empty"
	
		state.claimed = false
		state.buildingOwner = nil
		state.slatedForDemolition = false
		state.assignmentID = -1
		state.currentJobIndex = 1
		
		state.curConstructionAssignment = nil
		state.curDestructionAssignment = nil
		state.curAssignment = nil
		
		state.repair_material_dropped_off = false
		state.repair_materials = {}
		
		state.buildingUpgradeLevel = 0
		
		state.doors = 0
	>>

	receive odinBuildingCompleteMessage( int handle, gameSimJobInstanceHandle ji )
	<<
		send("gameSession","incSessionInt","buildingCount",1)
		--[[
		--Set up next suggested building here.
		--First check for emergencies. Food / military.
		local button1
		local button2
		local button3
		local foodCount = query("gameSession", "getSessionInt", "rawFoodCount")[1] + query("gameSession", "getSessionInt", "cookedFoodCount")[1]
		local foodFarms = query("gameSession","incSessionInt","foodFarmCount")[1]
		local dayCount = query("gameSession", "getSessionInt", "foodFarmCount")[1]
		local pop = query("gameSession", "getSessionInt", "lowerClassPopulation")[1] + query("gameSession", "getSessionInt", "middleClassPopulation")[1]
		if (foodCount < pop) or ((foodfarms < (math.floor(pop/25))) and (pop < 100)) then --You're likely starving OR you should probably have a farm anyway. After 100 pop I assume you can handle yourself tho
				--Next let's determine what your best food farm is.
				local biome = query("gameSession", "setSessionString", "biome")[1]
				if biome == "temperate" then
					if query("gameSession","getSessionBool","pumpkin_technology_unlocked")[1] == true then
						
					elseif query("gameSession","getSessionBool","wheat_technology_unlocked")[1] == true then
						
					elseif query("gameSession","getSessionBool","grape_technology_unlocked")[1] == true then
						
					else --corn
						
					end
				elseif biome == "tropical" then
					
				elseif biome == "desert" then
					
				else --cold
					query("gameSession","getSessionBool","wheat_technology_unlocked")[1]
				end
				
				
				send("odinRendererAdvisorMessage", "setNewAdvisedBuilding",
					button1, --button
					button2, --subbutton
					button3--subsubbutton if necessary
				    )
			end
		end
		--Prelumber Tier IE: BUILD A LUMBER
		if query("gameSession","getSessionBool","builtCarpentry Workshop")[1] == false then
			
		elseif --Lumber Tier
		
		elseif--Ceramics Tier
	
		else --Metalworks Tier.
			
		end
		]]
		-- update building stats here.
		local building_count = query("gameSession","getSessionInt","buildingCount")[1]
		if building_count > query("gameSession","getSessionInt","highestBuildingCount")[1] then	
			send("gameSession","setSessionInt","highestBuildingCount", building_count)
			send("gameSession","setSessionString","endGameString3", tostring(building_count) )
		end
          
          --Let's just track whether you have built stuff ever, too.
          send("gameSession","setSessionBool","built" .. state.buildingName, "true")
		
          local analytics = query("gameSession","getSessiongOH","analytics")
		if analytics and analytics[1] then
			send(analytics[1],"recheckAchievements")
		end
	
		state.rOH = SELF.id
		state.completed = true
	
		state.material_names = {}
		for i = 1, #state.materials do
			state.material_names[i] = query(state.materials[i],"getName")[1] -- cache name.
			send(state.materials[i], "DestroyedMessage")
			
			local resultROH = query( state.materials[i], "ROHQueryRequest" )
			send("rendStaticPropClassHandler", "odinRendererDeleteStaticProp", resultROH[1])
			destroyfromjob(state.materials[i], ji)
		end
		state.materials = {}
		
		state.contents = {}
		
		local dayCount = query("gameSession", "getSessionInt", "dayCount")[1]
		printl("analytics", "Player Built " .. state.buildingName .. " on Day " .. dayCount)
		-- end analytics.
		
		state.parent = EntityDB[state.buildingName]
		if state.parent == nil then
			printl("building", "WARNING failed to find edb information")
			return
		end

		if state.parent.tags then
			SELF.tags = state.parent.tags
		end
		
		SELF.tags.building_attack_target = true
          
		-- analytics stuff.
		if SELF.tags.workshop then
			send("gameSession", "incSessionInt", "workshopsBuilt", 1)
		end
		if SELF.tags.barracks then
			send("gameSession", "incSessionInt", "barracksBuilt", 1)
          end
		
		incMusic(3,30)
	>>

	receive InteractiveMessage( string messagereceived )
	<<
		printl("buildings", tostring(SELF.id) .. " Building Message Received: " .. messagereceived )
		
		--[[if messagereceived == "UpgradeBuilding" then
			SELF.tags["upgrade_in_progress"] = true
			state.curConstructionAssignment = query("gameSimulationCommandManager", "odinUpgradeBuildingMessage", SELF, #state.squares)
		elseif messagereceived == "CancelUpgradeBuilding" then
			if state.curConstructionAssignment ~= nil and SELF.tags["upgrade_in_progress"] == true then
				SELF.tags["upgrade_in_progress"] = nil
				-- cancel assignment
				state.curConstructionAssignment = nil
				send("rendUIManager", "SetOfficeBool", SELF, "upgrade_in_progress", false);
			end]]
			
		if messagereceived == "demolish" then
			printl("buildings", "slating building "  .. tostring(SELF.id) .. " for demolition...")
			send(SELF, "Demolition")
		end
	>>

	receive buildingSetConstructionAssignment ( gameSimAssignmentHandle assignment )
	<<
		state.curConstructionAssignment = assignment
		for i=1,#state.squares do		
			local results = query("gameSpatialDictionary",
							  "allObjectsInRadiusRequest",
							  state.squares[i], 1, true)
			send(results[1], "ClearNature", state.curConstructionAssignment)
		end
		send("rendBeaconClassHandler", "DeleteAssignmentBeacon", state.curConstructionAssignment);
	>>

	receive buildingSetDestructionAssignment ( gameSimAssignmentHandle assignment )
	<<
		state.curDestructionAssignment = assignment
	>>

	respond buildingGetConstructionAssignment()
	<<
		return "buildingSetConstructionAssignment", state.curConstructionAssignment
	>>
	
	receive odinBuildingClearSquaresMessage ()
	<<
		for k in pairs (state.squares) do
			state.squares[k] = nil
		end
	>>
	respond buildingGetFloorSize ()
	<<
		return "buildingGetFloorSize", #state.squares
	>>

	receive odinBuildingPointMessage ( int x, int y )
	<<
		gp = gameGridPosition:new()
		gp.x = x
		gp.y = y
		state.points[#state.points + 1] = gp
	>>

	receive odinBuildingSquareMessage (int x, int y, bool interior )
	<<
		gp = gameGridPosition:new()

		gp.x = x
		gp.y = y
		
		state.squares[#state.squares + 1] = gp

		send("gameSpatialDictionary", "gridSetCivilization", gp, 10)
		send("gameWorkshopManager", "SetNumInteriorSquares", SELF, #state.squares);
		
		-- HAX: send generic clear to nature objects	
		--[[local r = rand(1,10)
		if r == 1 then
			send("rendCommandManager",
				"odinRendererCreateParticleSystemMessage",
				"DustPuffMassive",
				gp.x,
				gp.y)
		end--]]
	>>
	
	receive addMaterialMessage(gameObjectHandle handle, string tag)
	<<
		-- make sure this handle stays LOCKED
		send("gameUtilitySingleton", "odinLockObject", handle);
		state.materials[#state.materials +1 ] = handle		
		send(handle, "BuildingLocked")
		send("rendOdinBuildingClassHandler", "odinRendererIncModelCompleteCosts", SELF.id, tag, 1)

	>>
	
	receive buildingAddRequiredModuleTagMessage(string name, int amount)
	<<

	>>

	receive buildingAddModuleMessage(string name, gameObjectHandle handle)
	<<
		if VALUE_STORE.showBuildingDebugConsole then printl ("building", "Building received module " .. name ) end
		send(handle, "SetParentBuilding", SELF);
		state.modules[#state.modules + 1] = {name, handle}
	>>

	respond getBuildingPosition ()
	<<
		-- Try to find an open building position
		local r = rand(1, #state.squares)
		return "getBuildingPositionReply", state.squares[r] 
	>>

	respond getBuildingGridPoint ()
	<<
		-- Try to find an open building position
		local r = rand(1, #state.points)

		return "getBuildingPositionReply", state.points[r] 
	>>

	respond buildingHasModules()
	<<
		local numModules = 0
		for i = 1, #state.modules do
			local result = query(state.modules[i][2], "ModuleIsPartOfBuilding");
			if result and result[1] == false then
				numModules = numModules + 1
			end
		end
		if numModules > 0 then
			return "buildingHasModulesResponse", true
		else
			return "buildingHasModulesResponse", false
		end
	>>

	respond buildingHasSquare(gameGridPosition square)
	<<
		for i = 1, #state.squares do
			if state.squares[i].x == square.x and state.squares[i].y == square.y then
				return "buildingSquareResponse", true
			end
		end
		
		return "buildingSquareResponse", false
	>>

	respond getBuildingModules( string name, gameObjectHandleVector list )
	<<
		for i = 1, #state.modules do
			if name:len() == 0 or state.modules[i][1] == name then
				list[#list + 1] = state.modules[i][2]
			end
		end

		return "buildingModules", list
	>>

	respond gridGetPosition()
	<<
		return "reportedPosition", state.squares[1]
	>>

	respond GetRandomBuildingPosition()
	<<			
		return "reportedPosition", state.squares[rand(1,#state.squares)]
	>>

	respond gridReportPosition()
	<<			
		return "reportedPosition", state.squares[1]
	>>

	respond isBuildingComplete()
	<<
		return "buildingComplete", state.completed
	>>

	respond isBuildingClaimed()
	<<
		return "buildingClaimed", state.claimed
	>>
	
	receive DestroyBuilding( gameSimJobInstanceHandle ji )
	<<
		if SELF.tags.pub then
			-- dump the booze!
			function dropBooze( v )
				local results = query("scriptManager",
							"scriptCreateGameObjectRequest",
							"item",
							{ legacyString = v} )[1]
							
				send(results,"ClaimItem")
				
				local positionResult = query(SELF, "GetRandomBuildingPosition")[1]
				local x = positionResult.x
				local y = positionResult.y
			
				-- drop outside foundation so we don't get floaters.
				local isInvalidDrop = true
				local i = 0
				while isInvalidDrop do
					if i > 0 then
						positionResult.x = x + rand(i * -1,i)
						positionResult.y = y + rand(i * -1,i)
					end
					isInvalidDrop = query("gameSpatialDictionary","gridHasSpatialTag",positionResult,"occupiedByStructure" )[1]
					i = i + 1
				end
				
				send(results,"GameObjectPlace",positionResult.x,positionResult.y  )
			end
			if state.brewTable and #state.brewTable > 0 then
				for k,v in pairs(state.brewTable) do
					dropBooze( v )	
				end
			end
			if state.spiritsTable and #state.spiritsTable > 0 then
				for k,v in pairs(state.spiritsTable) do
					dropBooze( v )	
				end
			end
			if state.laudanumTable and  #state.laudanumTable > 0 then
				for k,v in pairs(state.laudanumTable) do
					dropBooze( v )	
				end
			end
		end
		
		send("gameSession","incSessionInt","buildingCount",-1)
		
		for i = 1, #state.modules do
			send(state.modules[i][2], "DestroyModule", ji);
		end

		-- return materials.
		if state.material_names then 
			for k,v in pairs(state.material_names) do
				
				local results = query("scriptManager",
										"scriptCreateGameObjectRequest",
										"item",
										{ legacyString = v} )[1]
						
				send(results,"ClaimItem")
				
				local positionResult = query(SELF, "GetRandomBuildingPosition")[1]
				
				local x = positionResult.x
				local y = positionResult.y
			
				-- drop outside foundation so we don't get floaters.
				local isInvalidDrop = true
				local i = 0
				while isInvalidDrop do
					if i > 0 then
						positionResult.x = x + rand(i * -1,i)
						positionResult.y = y + rand(i * -1,i)
					end
					isInvalidDrop = query( "gameSpatialDictionary",
									"gridHasSpatialTag",
									positionResult,
									"occupiedByStructure" )[1]
	
					i = i + 1
				end
				
				send(results,"GameObjectPlace",positionResult.x,positionResult.y  )
			
				send("rendCommandManager",
					"odinRendererCreateParticleSystemMessage",
					"DustPuffV1",x,y)
			end
		end

		send("gameWorkshopManager", "RemoveOffice", SELF);
		send("gameWorkshopManager", "RemoveWorkshop", SELF);

		send("rendOdinBuildingClassHandler", "odinRendererDeleteBuildingMessage", SELF);
		send("gameSpatialDictionary", "gridRemoveObject", SELF);
		send("gameSpatialDictionary", "removeBuilding", SELF);
		if ji == nil then
			destroy(SELF);
		else
			destroyfromjob(SELF, ji)
		end
	>>

	receive Demolition()
	<<
		if state.buildingOwner ~= nil then
			-- remove old owner	
			local owner = state.buildingOwner
			send(state.buildingOwner, "claimWorkBuilding", nil)
		end
	
		if state.slatedForDemolition == true then
			-- Already received demolition
			return
		end
		
		SELF.tags.slated_for_demolition = true
		if state.curConstructionAssignment ~= nil then
			send("gameBlackboard", "cancelAssignment", state.curConstructionAssignment);
			state.curConstructionAssignment = nil
		end

		if SELF.tags.office then
			
			-- if trade office and there's a foreign trade group, cancel their mission.
			if SELF.tags.trade_office then
				printl("buildings", "dismantling trade office: sending cancel mission to all active traders")
				--send( query("gameSession","getSessiongOH", "Empire")[1], "cancelTradeMissions" )
				send( query("gameSession","getSessiongOH", "Stahlmark")[1], "cancelTradeMissions" )
				send( query("gameSession","getSessiongOH", "Republique")[1], "cancelTradeMissions" )
				send( query("gameSession","getSessiongOH", "Novorus")[1], "cancelTradeMissions" )
			end
			
			send("rendUIManager", "CloseProductionMenu", SELF)
			send(SELF, "setBuildingOwner", nil)
			local newJobs = {}
			state.jobs = newJobs
			state.slatedForDemolition = true
			
		elseif SELF.tags.workshop then
			
			state.slatedForDemolition = true
			send("rendUIManager", "CloseProductionMenu", SELF)
			send("gameWorkshopManager", "RemoveAllWorkshopJobs", SELF);
			send(SELF, "setBuildingOwner", nil)
			local newJobs = {}
			state.jobs = newJobs
			
		elseif SELF.tags.house then
			if SELF.tags.lower_class_house then
				send("gameSession", "incSessionInt", "LcPopulationAllowed", state.lc_pop_cap_increase * -1)
				send("gameSession", "incSessionInt", "totalPopulationAllowed", state.lc_pop_cap_increase * -1)
			elseif SELF.tags.middle_class_house then
				send("gameSession", "incSessionInt", "McPopulationAllowed", state.mc_pop_cap_increase * -1)
				send("gameSession", "incSessionInt", "totalPopulationAllowed", state.mc_pop_cap_increase * -1)
			end
			state.slatedForDemolition = true
			send("rendUIManager", "CloseProductionMenu", SELF)
			send("rendUIManager", "CloseHousingMenu", SELF)
		end
		
		local assignmentResults = query("gameBlackboard",
								  "gameObjectNewAssignmentMessage",
								  SELF,
								  "Demolish Building",
								  "construction",
								  "construction")
		
		state.curDestructionAssignment = assignmentResults[1]

		send("gameBlackboard",
			"gameObjectNewJobToAssignment",
			assignmentResults[1],
			SELF,
			"Disassemble Building",
			"building",
			true )
		
		--send("gameWorkshopManager", "RemoveOffice", SELF);
		--send("gameWorkshopManager", "RemoveWorkshop", SELF);
	>>

	respond buildingHasNoCompletedModules()
	<<
		local list = {}
		for i = 1, #state.modules do
			local moduleComplete = query(state.modules[i][2], "isModuleComplete")
			if moduleComplete and moduleComplete[1] == true then
				local moduleCountsAsPartOfBuilding = query(state.modules[i][2], "ModuleIsPartOfBuilding");
				if moduleCountsAsPartOfBuilding and not moduleCountsAsPartOfBuilding[1] then
					return "buildingModuleCompletionState", false
				end
			end
		end
		return "buildingModuleCompletionState", true
	>>

	receive buildingClaimed ( bool claimState )
	<<
		if claimState then
			if VALUE_STORE["showBuildingDebugConsole"] then printl("buildings", "Building Claimed: TRUE") end
		else
			if VALUE_STORE["showBuildingDebugConsole"] then printl("buildings", "Building Claimed: FALSE") end
		end
	>>

	respond getBuildingOwner()
	<<
		return "buildingOwner", state.buildingOwner
	>>

	receive setBuildingOwner(gameObjectHandle newOwner)
	<<
		if state.slatedForDemolition then
			-- do alert here?
			return
		end
		
		if state.buildingOwner ~= nil and state.buildingOwner ~= newOwner then
			-- remove from old owner	
			local owner = state.buildingOwner
			send(state.buildingOwner, "claimWorkBuilding", nil)
			
			-- change class of existing owner/workparty if necessary
			--[[if SELF.tags.office then
				local workPartyResults = query("gameBlackboard",
					"GetWorkPartyWorkers",
					owner)[1]
				
				send(workPartyResults, "RevertOfficeWorker")
				--send("gameBlackboard","RevertOfficeGroup",owner)
			end]]
		end
	
		state.buildingOwner = newOwner
		if not state.buildingOwner then
			state.claimed = false
			SELF.tags.overseer_active = nil
			
			if state.eventLocked == true then
				-- does this fail the event? let the director figure it out
				send(state.event_director,"setKeyBool","eventBuildingOwnerNil",true)
				send(state.event_director,"releaseEventBuilding")
			end
		else
			state.claimed = true
			--local workPartyResults = query("gameBlackboard",
			--		"GetWorkPartyWorkers",
			--		newOwner)[1]
			
			send(newOwner, "claimWorkBuilding", SELF)
			
			--[[if SELF.tags.office then
				send(workPartyResults,
					"BecomeOfficeWorkers",
					string.gsub( string.gsub( string.lower(state.buildingName), " ", "_" ), "'", "")
					)
			end]]
		end
	>>

	respond sleepInfo()
	<<
		-- This is the mood change for sleeping on the floor of a building
		return "sleepInfoResponse", "buildingFloor"
	>>

	respond getSleepDirection()
	<<
		return "sleepDirectionResponse", ""
	>>

	receive AssignmentSuspendedMessage ( gameSimAssignmentHandle assignment )
	<<
		if state.curConstructionAssignment == assignment then
			return
		end
	>>

	receive AssignmentCancelledMessage( gameSimAssignmentHandle assignment )
	<<
		if state.curConstructionAssignment == assignment then
			state.curConstructionAssignment = nil

			local newModules = {}
			local deleteModules = {}

			for i = 1, #state.modules do
				local moduleComplete = query(state.modules[i][2], "isModuleComplete");
				if moduleComplete and moduleComplete[1] == false then
					deleteModules[#deleteModules + 1] = state.modules[i]
				else
					newModules[#newModules + 1] = state.modules[i]
				end
			end

			for i = 1, #deleteModules do				
				if VALUE_STORE["showBuildingDebugConsole"] then printl("building", "deleting module " .. deleteModules[i][1]) end
				send(deleteModules[i][2], "BuildingDeleted");
			end

			state.modules = newModules

			if not state.completed then
				send("rendOdinBuildingClassHandler", "odinRendererDeleteBuildingMessage", SELF);
				send("gameSpatialDictionary", "gridRemoveObject", SELF);
				send("gameSpatialDictionary", "removeBuilding", SELF);
				for i = 1, #state.materials do
					send("gameUtilitySingleton", "odinUnlockObject", state.materials[i])
					send(state.materials[i], "BuildingUnlocked");
				end
				send("gameWorkshopManager", "RemoveOffice", SELF);
				send("gameWorkshopManager", "RemoveWorkshop", SELF);
				destroy(SELF)
			end
			return
		end
	>>
	
     respond getTags()
     <<
          return "getTagsResponse", SELF.tags
     >>
	
	receive addTag( string name )
	<<
		SELF.tags[name] = true
	>>
	
	receive removeTag( string name )
	<<
		SELF.tags[name] = nil
	>>
	
	respond getStairStyle()
	<<
		if state.parent and state.parent.stairStyle then
			return "getStairStyleMessage", state.parent.stairStyle
		end
		return "getStairStyleMessage", "wood"
	>>
	
	respond getBuildingQuality()
	<<
		return "buildingQuality", state.buildingQuality, state.buildingQualityName
	>>
	
	receive BuildingContainerAddItem( gameObjectHandle goh)
	<<
		state.contents[#state.contents + 1] = goh
	>>
	
	respond BuildingContainerRemoveItem()
	<<
		if #state.contents == 0 then
			return "BuildingContainerRemoveItemMessage", false
		else
			local x = state.contents[#state.contents]
			table.remove(state.contents)
			destroy(x)
			return "BuildingContainerRemoveItemMessage", true
		end
	>>
	
	respond getContainerAmount()
	<<
		return "getContainerAmountResponse", #state.contents
	>>
	
	receive damageMessage( gameObjectHandle attacker, string damageType, int damageAmount, string onhit_effect )
	<<
		-- Do a point of damage here.
		state.hitpoints = state.hitpoints - 1
		send("rendUIManager", "SetOfficeInt", SELF, "buildingHP", state.hitpoints)

		local buildingName = "a " .. state.buildingName
		if state.buildingFancyName ~= nil and state.buildingFancyName ~= "" then
			buildingName = state.buildingFancyName
		end
		
		local attackerTags = query(attacker, "getTags")[1]
		local tagsToDetect = {"fishperson", "citizen", "bandit", "foreigner"}
		local identifier = "mysterious vandal"
		if attacker == SELF.id then
			identifier = " wear and tear"
		end
		for k,v in pairs(attackerTags) do
			for k2=1, #tagsToDetect do
				if k == tagsToDetect[k2] then
					identifier = tagsToDetect[k2]
					break
				end
			end
		end
		
		-- just for fun
		local adjectives = { "disrespectful", "vandalous", "terrible", "callous", "rampaging", "dangerous", "dastardly", "", "", "", "violent"}
		identifier = adjectives[ rand(1,#adjectives)] .. " " .. identifier
		
		if state.hitpoints > 0 then
			local damageState = "ERROR"
			if (state.hitpoints < state.hitpointsmax) and (state.hitpoints <= 2) then
				damageState = "nearly destroyed"
			elseif state.hitpoints < divFloor(state.hitpointsmax,2) then --half HP or less
				damageState = "severely damaged"
			elseif state.hitpoints < state.hitpointsmax then
				damageState = "somewhat damaged"
			else --...somehow you have full HP???
				damageState = "undamaged, somehow"
			end
			
			send("rendUIManager", "SetOfficeString", SELF, "buildingHPDescription", damageState)
			
			send("rendCommandManager",
				"odinRendererStubMessage",
				"ui\\thoughtIcons.xml", -- iconskin
				"explosion", -- icon
				"Structure under attack!", -- header text
				"A " .. identifier .. " is attacking " .. buildingName .. ", leaving it " .. damageState .. ".", -- text description
				"Left-click to zoom. Right-click to dismiss.", -- action string
				"building_damage_" .. state.buildingName, -- alert type (for stacking)
				"ui//eventart//modules_breaking.png", -- imagename for bg
				"critical", -- importance: low / high / critical
				attacker, -- object ID
				60 * 1000, -- duration in ms
				30 * 1000, -- "snooze" time if triggered multiple times in rapid succession
				nil) -- gameobjecthandle of director, null if none
			
			-- now make a repair job!
			if not state.curRepairAssignment then
				state.curRepairAssignment = query("gameBlackboard",
											"gameObjectNewAssignmentMessage",
											SELF,
											"Repair Building",
											"",
											"")[1]
			end
			
			send("gameBlackboard",
				"gameObjectNewJobToAssignment",
				state.curRepairAssignment,
				SELF,
				"Repair Building",
				"building",
				true )
			
		else
			--boooom

			send("rendCommandManager",
				"odinRendererFYIMessage",
				"ui\\thoughtIcons.xml",													-- iconskin
				"explosion",															-- icon
				state.buildingName .. " destroyed!",									-- header text
				"A " .. identifier .. " has destroyed " .. buildingName .. ".",	-- text description
				"Left-click to zoom. Right-click to dismiss.", -- action string
				"building_damage_" .. state.buildingName,								-- alert type (for stacking)
				"ui//eventart//modules_breaking.png",									-- imagename for bg
				"critical",																-- importance: low / high / critical
				attacker,																-- object ID
				60 * 1000,																	-- duration in ms
				0, -- "snooze" time if triggered multiple times in rapid succession
				nil) -- gameobjecthandle of director, null if none
			
			-- more rubble/debris, less explosions.
			
			for k,v in pairs(state.squares) do
				local r = rand(1,10)
				
				if r > 6 then
					send("rendCommandManager",
						"odinRendererCreateParticleSystemMessage",
						"DustPuffMassive",
						v.x,
						v.y)
				end
				
				if r > 5 then
					local debrisNames = EntityDB.WorldStats.buildingDebrisList
					local debrisName = debrisNames[ rand(1,#debrisNames)]
					
					local handle = query( "scriptManager",
										"scriptCreateGameObjectRequest",
										"clearable",
										{legacyString= debrisName })[1]
					
					send(handle, "GameObjectPlace", v.x, v.y)
					
				elseif r==1 then
					
					local explosionNames = {
						"Landmine Explosion",
						"Medium Explosion",
						"Ammo Explosion",
					}
				
					local boomName = explosionNames[ rand(1,#explosionNames)]
					local handle = query( "scriptManager",
										"scriptCreateGameObjectRequest",
										"explosion",
										{legacyString= boomName })[1]
					
					send(handle, "GameObjectPlace", v.x, v.y)
				end
			end
			
			-- unlock construction materials if unfinished
			if not state.completed then
				for i = 1, #state.materials do
					send("gameUtilitySingleton",
						"odinUnlockObject",
						state.materials[i])
					
					send(state.materials[i], "BuildingUnlocked")
				end
			end

			if SELF.tags.lower_class_house then
				send("gameSession", "incSessionInt", "LcPopulationAllowed", -2)
				send("gameSession", "incSessionInt", "totalPopulationAllowed", -2)
			elseif SELF.tags.middle_class_house then
				send("gameSession", "incSessionInt", "McPopulationAllowed", -1)
				send("gameSession", "incSessionInt", "totalPopulationAllowed", -1)
			end
			
			-- and kill it
			send(SELF, "setBuildingOwner", nil)
			send(SELF,"DestroyBuilding",nil)
		end
	>>
	
	receive BuildingRepairMessage(gameSimJobInstanceHandle ji)
	<<
		state.hitpoints = state.hitpoints + 1
		send("rendUIManager", "SetOfficeInt", SELF, "buildingHP", state.hitpoints)

		if state.hitpoints < divFloor(state.hitpointsmax,2) then --half HP or less
			send("rendUIManager", "SetOfficeString", SELF, "buildingHPDescription", "Very Damaged")
		elseif state.hitpoints < state.hitpointsmax then
			send("rendUIManager", "SetOfficeString", SELF, "buildingHPDescription", "Damaged")
		else
			send("rendUIManager", "SetOfficeString", SELF, "buildingHPDescription", "Undamaged")
			state.curRepairAssignment = nil
		end
	>>
	
	receive lockEventBuilding(gameObjectHandle director)
	<<
		printl("events", " received lockEventBuilding from event_director: " .. query(director,"getName")[1] )
		state.eventLocked = true
		SELF.tags.eventlocked = true
		state.event_director = director
		
		local arcName = query(state.event_director,"getEventArcName")[1]
		local eventBuildingArt = query(state.event_director,"getEventBuildingArt")[1]
		local eventBuildingText = query(state.event_director,"getEventBuildingText")[1]
		
		if SELF.tags.laboratory then
			-- if building is a lab, make sure we have a chalkboard. If not, give warning.
			local modules = query("gameWorkshopManager", "getBuildingModulesGameSide", SELF)[1]
			local found_chalkboard = false
			for k,v in pairs(modules) do
				if v.tags.lab_equipment then
					found_chalkboard = true
				end
			end
			if not found_chalkboard then
				-- you need a chalkboard!
				
				send("rendCommandManager",
					"odinRendererStubMessage", -- "odinRendererStubMessage",
					"ui\\orderIcons2.xml", -- iconskin
					"chalkboard", -- icon
					"Laboratory Needs Equipment", -- header text
					"The Laboratory cannot engage in a 'special research project' without lab equipment. Order a lab modules constructed!", -- text description
					"Left-click for more information. Right-click to dismiss.", -- action string
					"buildingMissingModuleAlert", -- alert type (for stacking)
					"ui\\eventart\\doing_science.png", -- imagename for bg
					"low", -- importance: low / high / critical
					SELF.id, -- object ID
					60 * 1000, -- duration in ms
					0, -- snooze
					state.event_director) 
				
			end
		elseif SELF.tags.barracks then
			
			-- set soldiers to interrogate people.
			-- should this alert go here, or is it handled by events?
			send("rendCommandManager",
				"odinRendererStubMessage", -- "odinRendererStubMessage",
				"ui\\orderIcons.xml", -- iconskin
				"barracks", -- icon
				"Soldiers To Perform Interrogations", -- header text
				"Solders will perform interogations of random colonists. The truth will be found regardless of hurt feelings!", -- text description
				"Left-click for more information. Right-click to dismiss.", -- action string
				"buildingMissingModuleAlert", -- alert type (for stacking)
				"ui\\eventart\\squad.png", -- imagename for bg
				"low", -- importance: low / high / critical
				SELF.id, -- object ID
				60 * 1000, -- duration in ms
				0, -- snooze
				state.event_director) 
			
		elseif SELF.tags.foreign_office then
			-- if building is a lab, make sure we have a chalkboard. If not, give warning.
			local modules = query("gameWorkshopManager", "getBuildingModulesGameSide", SELF)[1]
			local found_desk = false
			for k,v in pairs(modules) do
				if v.tags.diplomatic_desk then
					found_desk = true
				end
			end
			if not found_desk then
				-- you need a chalkboard!
				
				send("rendCommandManager",
					"odinRendererStubMessage", -- "odinRendererStubMessage",
					"ui\\orderIcons.xml", -- iconskin
					"standing_desk_icon", -- icon
					"Foreign Office Needs Desks", -- header text
					"The Foreign Office cannot engage in a 'special research project' without working desks. Order standing desks constructed!", -- text description
					"Left-click for more information. Right-click to dismiss.", -- action string
					"buildingMissingModuleAlert", -- alert type (for stacking)
					"ui\\eventart\\doing_science.png", -- imagename for bg
					"low", -- importance: low / high / critical
					SELF.id, -- object ID
					60 * 1000, -- duration in ms
					0, -- snooze
					state.event_director) 
				
			end
			
		elseif SELF.tags.chapel then
			-- if building is a chapel, boot out anyone waiting for confession
			-- ping from random square in chapel
			
			local r = rand(1, #state.squares)
			local results = query("gameSpatialDictionary",
								"allObjectsInRadiusWithTagRequest",
								state.squares[r],
								10,
								"citizen",
								true)
			
			if results and results[1] then
				send(results[1], "removeTag", "waiting_for_confession")
			end
		end
		
		send(SELF,"setEventBuildingTitleArtText", arcName, eventBuildingArt, eventBuildingText)
	>>
	
	receive releaseEventBuilding()
	<<
		if state.event_director then 
			printl("events", " received releaseEventBuilding from event_director: " .. query(state.event_director,"getName")[1] )
			state.eventLocked = false
			SELF.tags.eventlocked = nil
			state.event_director = nil
			
			send("rendUIManager", "SetOfficeString", SELF, "eventName", "")
			send("rendUIManager", "SetOfficeString", SELF, "eventArt", "")
			send("rendUIManager", "SetOfficeString", SELF, "eventText", "")
		end
	>>
	
	receive earnEventPoint( gameObjectHandle worker )
	<<
		printl("events", "building received earnEventPoint")
		if state.event_director then 
			send(state.event_director,"incKeyInt","eventPoints",1)
		else
			printl("buildings", "Warning! received earnEventPoint but has no registered event_director")
		end
	>>
	
	receive setEventBuildingTitleArtText(string title, string art, string text)
	<<
		if title then
			send("rendUIManager", "SetOfficeString", SELF, "eventName", title)
		end
		if art then
			send("rendUIManager", "SetOfficeString", SELF, "eventArt", art)
		end
		if text then
			send("rendUIManager", "SetOfficeString", SELF, "eventText", text)
		end
	>>
	
	respond getRandomModuleByTag(string tag)
	<<
		local modules = query("gameWorkshopManager", "getBuildingModulesGameSide", SELF)[1]
		local module = false
		for k,v in pairs(modules) do
			local moduleTags = query(k, "getTags")[1]
			if moduleTags[tag] then
				-- found one.
				module = v
			end
		end
		
		if module == false then
			return "getRandomModuleByTagResponse", nil
		end
		
		return "getRandomModuleByTagResponse", module
	>>
	
	respond checkForModuleTagInBuilding( string tag)
	<<
		local response = nil
		local modules = query("gameWorkshopManager", "getBuildingModulesGameSide", SELF)[1]
		for k,v in pairs(modules) do
			local tags = query(v,"getTags")[1]
			if tags[tag] then
				response = true
			end
		end
		return "checkForModuleTag", response
	>>
	
	receive recalculateQuality()
	<<

		if state.doors and state.doors < 1 then
               send("rendCommandManager",
				"odinRendererStubMessage",
				"ui\\orderIcons.xml", -- iconskin
				"lower_class_door_icon", -- icon
				"Building has no Door!", -- header text
				"Your newly constructed " .. state.buildingName .. " has no door! Any colonists inside will be trapped until you add a door to the building.", -- text description
				"Left-click to zoom, right click to dismiss.", -- action string
				"above_100chars", -- alert type (for stacking)
				"ui\\eventart\\immigration.png", -- imagename for bg
				"high", -- importance: low / high / critical
				state.rOH, -- object ID
				60 * 1000, -- duration in ms
				0, -- "snooze" time if triggered multiple times in rapid succession
				nil)
               
               send("rendCommandManager",
				"odinRendererPlaySoundMessage",
				"alertBad")
          end
		
		local newQuality = 1 -- defaults to 1 now! Yay!
		local modules = query("gameWorkshopManager", "getBuildingModulesGameSide", SELF)[1]

		local doTutorial = false
		if query("gameSession", "getSessionBool", "caseTutorialActive")[1] == true then
			doTutorial = true
		end
          
		for k,v in pairs(modules) do
			local moduleName = query(k, "getModuleName")[1]
			
			if doTutorial == true then
				if SELF.tags.workshop then
					if (moduleName == "Carpentry Workbench") or
						(moduleName == "Small Stone Oven") then
						
						local director = query("gameSession","getSessiongOH","event_director_tutorial")[1]
						send(director, "setKeyString", "moduleBench", "yes")
					end
				elseif SELF.tags.house then
					if (moduleName == "Cot") or
						(moduleName == "Lower Class Bed") or
						(moduleName == "Middle Class Bed") or
						(moduleName == "Upper Class Bed") then
						
						local director = query("gameSession","getSessiongOH","event_director_tutorial")[1]
						send("gameSession", "setSessionBool", "bedTutorialDone", true)
						send(director, "setKeyString", "moduleBed", "yes")
					end
				end
			end
			
			local entityData = EntityDB[moduleName]
			if entityData.quality then
				newQuality = newQuality + entityData.quality
			end
		end
		
		state.buildingQuality = newQuality
		
		if state.buildingQuality >= 6 then -- 6+ is amazing
			state.buildingQualityName = "Excellent"
			if doTutorial then
				local director = query("gameSession","getSessiongOH","event_director_tutorial")[1]
				send(director, "setKeyString", "qualityDone", "yes")
			end
		elseif state.buildingQuality >= 4 then -- 4, 5 is great
			state.buildingQualityName = "Great"
			if doTutorial then
				local director = query("gameSession","getSessiongOH","event_director_tutorial")[1]
				send(director, "setKeyString", "qualityDone", "yes")
			end
		elseif state.buildingQuality >= 2 then -- 2, 3 is good
			state.buildingQualityName = "Comfortable"
		elseif state.buildingQuality >= -1 then -- 1, 0, or -1 is normal
			state.buildingQualityName = "Typical"
		elseif state.buildingQuality >= -3 then -- -2, -3 is bad
			state.buildingQualityName = "Uncomfortable"
		elseif state.buildingQuality >= -5 then -- -4, -5 is really bad.
			state.buildingQualityName = "Distasteful"
		else -- -6 or worse is THE WORSTEST
			state.buildingQualityName = "Horrible"
		end
		
		send("rendUIManager","SetOfficeInt", SELF, "buildingQuality", state.buildingQuality) -- used for tooltips
		
		if state.buildingQuality > 0 then
			send("rendUIManager","SetOfficeString", SELF, "buildingQualityBad", "")
			send("rendUIManager","SetOfficeString", SELF, "buildingQualityNeutral", "")
			send("rendUIManager","SetOfficeString", SELF, "buildingQualityGood", "+".. tostring(state.buildingQuality))
			
		elseif state.buildingQuality == 0 then
			send("rendUIManager","SetOfficeString", SELF, "buildingQualityBad", "")
			send("rendUIManager","SetOfficeString", SELF, "buildingQualityNeutral", "0")
			send("rendUIManager","SetOfficeString", SELF, "buildingQualityGood", "")
			
		else
			send("rendUIManager","SetOfficeString", SELF, "buildingQualityBad", tostring(state.buildingQuality))
			send("rendUIManager","SetOfficeString", SELF, "buildingQualityNeutral", "")
			send("rendUIManager","SetOfficeString", SELF, "buildingQualityGood", "")
		end
		
		send("rendUIManager","SetOfficeString", SELF, "buildingQualityDescription", state.buildingQualityName)
	>>
	
	receive newShiftUpdate()
	<<

	>>
	
	respond GetBuildingFancyName()
	<<
		if not state.buildingFancyName or
			state.buildingFancyName == "" then
			
			return "buildingFancyNameResult", state.buildingName
		end
		return "buildingFancyNameResult", state.buildingFancyName
	>>
	
	respond getBuildingName()
	<<
		return "buildingNameRespond", state.buildingName
	>>
	
     receive addDoor()
	<<
          state.doors = state.doors + 1
	>>
     
     receive removeDoor()
	<<
          state.doors = state.doors - 1
	>>
     
	respond hasDoor()
	<<
		if not state.doors or state.doors == 0 then
			return "hasDoorRespond", false
		end
		return "hasDoorRespond", true
	>>
	
     respond getName()
	<<
		return "buildingNameRespond", state.buildingName
	>>
	
	respond doesBuildingNeedSupplies( string g)
	<<
		local response = false
		
		-- what supplies are we being given?
		local supply_tags = {}
		for k,v in pairs(EntityDB[g].tags) do
			supply_tags[v] = true		
		end
		
		local my_supply_info = EntityDB[ state.buildingName ].supply_info
		for k,v in pairs(my_supply_info) do
			if SELF.tags["needs_resupply" .. tostring(v.tier ) ] and supply_tags[k] then
				response = true
				break
			end
		end
		
		-- then: do we need this type of supply? If yes, return true.
		return "doesBuildingNeedSuppliesResponse", response
	>>
	
	respond getRenderHandle()
	<<
		return "renderHandleResponse", state.rOH
	>>
>>