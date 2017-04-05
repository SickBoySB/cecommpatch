gameobject "module" inherit "spatialobject" -- inherit "object_damage" not yet - dgb
<<
	local 
	<<
	>>

	state
	<<
		gameObjectHandle parentBuilding
		int moduleRotation
		int moduleWidth
		int moduleHeight
		int local_player
		int local_id
		gameGridPosition machinePos
		string statedBuildingName
		string statedModuleName
		string processingJobName
		string moduleGhost
		string name
		
		bool processing
		bool repair_material_dropped_off
		string processingName
		string processingJobName
		int processingTime
		int processingJobID
		int renderHandle
		
		int timesUsed
		int usesUntilDamaged
		
		bool complete
		table materials
		gameObjectHandle owner

	>>

	receive BuildingDeleted()
	<<
		send(SELF,"DestroyModule", nil)
	>>
	
	receive DestroyModule(gameSimJobInstanceHandle ji)
	<<
		printl("buildings", "module received DestroyModule")
		send("rendGhostModelClassHandler","odinRendererDeleteGhostModel",SELF.id,true)
		
 		if state.complete then
			send("rendCommandManager",
				"odinRendererDeleteParticleSystemMessage",
				"BlueObject",
				state.position.x,
				state.position.y)
			
			
			send("rendMachineClassHandler","odinRendererDeleteMachineMessage",SELF.id)
			
			send("rendInteractiveObjectClassHandler",
				"odinRendererPlaySFXOnInteractive",
				SELF.id,
				"Tree Fall Hit")
			
			if state.material_names and
				not SELF.tags.needs_upkeep then
				
				for k,v in pairs(state.material_names) do
						
					local results = query("scriptManager",
									"scriptCreateGameObjectRequest",
									"item",
									{ legacyString = v} )[1]
					
					send(results,"ClaimItem")
					
					local positionResult = query("gameSpatialDictionary",
								   "nearbyEmptyGridSquare",
								   state.position,
									2)[1]
				
					local x = positionResult.x --+ rand(-1,1)
					local y = positionResult.y --+ rand(-1,1)
					
					send(results,"GameObjectPlace",positionResult.x,positionResult.y  )
				
					send("rendCommandManager",
						"odinRendererCreateParticleSystemMessage",
						"DustPuffV1",
						x,
						y)
				end
			end
			
			if SELF.tags.needs_upkeep then
				send("rendCommandManager",
					"odinRendererStubMessage",
					"ui\\thoughtIcons.xml", -- iconskin
					"hammer", -- icon
					state.statedModuleName .. " destroyed by dismantling", -- header text
					"A " .. state.statedModuleName .. " was ordered dismantled while still in need of repairs. \z
						It's just a bunch of useless, broken junk now and can't be re-installed.", -- text description
					"Right-click to dismiss.", -- action string
					"moduleDismantleUpkeepWarning", -- alert type (for stacking)
					"ui//eventart//modules_breaking.png", -- imagename for bg
					"critical", -- importance: low / high / critical
					nil, -- object ID
					45 * 1000, -- duration in ms
					0, -- "snooze" time if triggered multiple times in rapid succession
					nil ) -- gameobjecthandle of director, null if none
			end
		end
		

		local myTable = EntityDB[state.name]
		if myTable.standalone then
			send("gameSpatialDictionary", "RemoveStandaloneModule", SELF);
		end

		send("gameSpatialDictionary", "gridRemoveObject", SELF);
		destroyfromjob(SELF, ji);
	>>

	receive ModuleSetCosts(string name, string icon, int value)
	<<
		send("rendGhostModelClassHandler", "odinRendererSetModelCosts", SELF.id, name, icon, value);		
	>>

	receive addMaterialMessage(gameObjectHandle handle, string tag)
	<<	
		-- make sure this handle stays LOCKED
		send("gameUtilitySingleton", "odinLockObject", handle)
		state.materials[#state.materials + 1 ] = handle		
		send(handle, "BuildingLocked");
		send("rendGhostModelClassHandler", "odinRendererIncModelCompleteCosts", SELF.id, tag, 1)
	>>

	respond ModuleIsPartOfBuilding()
	<<
		return "moduleBuildingStatus", false				-- for now...
	>>

	respond isModuleComplete()
	<<
		return "moduleComplete", state.complete
	>>

	receive Create( stringstringMapHandle init )
	<<
		state.name = init["legacyString"]
		state.moduleWidth = tonumber(init["width"])
		state.moduleHeight = tonumber(init["height"])
		state.local_player = tonumber(init["local_player"])
		state.local_id = tonumber(init["local_id"])
		state.processing = false
		state.processingJobID = -1
		state.processingName = ""
		state.processingTime = 0
		state.timesUsed = 0
		state.repair_material_dropped_off = false
		state.owner = nil
		state.parentBuilding = nil
		
		local myTable = EntityDB[state.name]
		if myTable == nil then
			if VALUE_STORE["showModuleDebugConsole"] then printl("module", "module create: " .. state.name .. " not found...") end
 			scripterror("module creation: " .. state.name .. " not found in EDB")
			return
		end

		if myTable.tags then
			SELF.tags = myTable.tags
		end

		if myTable.standalone then
			send("gameSpatialDictionary", "AddStandaloneModule", SELF);
		end
		
		SELF.tags.under_construction = true
		
		-- TODO: allow these to be set PER MODULE if a field is set in modules.edb
		local worldstats = EntityDB["WorldStats"]
		state.usesUntilDamaged = worldstats["baseModuleUsesBeforeDamage"]

		SELF.tags["module"] = true
		SELF.tags["unclaimed_module"] = true
		state.complete = false
		
		ready()
	>>

	receive Update()
	<<
		-- see if we're still reachable
		local validresults = query("gameSpatialDictionary","hasValidAccessPoints",SELF)

		if validresults[ 1 ] == false then 
			if SELF.tags["unreachable"] == false then
				
				SELF.tags["unreachable"] = true
				printl("buildings","WARNING: setting module " .. state.name .. " to unreachable");
				if state.complete then
					send("rendMachineClassHandler", "odinRendererMachineSetBrokenMessage", state.renderHandle, true);
				end
			end
		else 
			if SELF.tags["unreachable"] == true then
				SELF.tags["unreachable"] = false
				printl("setting module " .. state.name .. " to reachable");
				if state.complete then
					send("rendMachineClassHandler", "odinRendererMachineSetBrokenMessage", state.renderHandle, false);
				end
			end
		end
	>>

	respond getAccessPosition()
	<<
		-- Return a position we can be accessed from.
		if not SELF.tags["front_access_module"] then
			return "moduleAccessPosition", state.position
		end
		-- frontaccessmodule returns a valid position
	>>

	receive ModuleSetCreationInformation(string moduleName, string buildingName, int modulePosX, int modulePosY, int moduleRotation, gameObjectHandle bh, bool immediate, string ghostModule)
	<<
		state.machinePos.x = modulePosX
		state.machinePos.y = modulePosY
		state.moduleRotation = moduleRotation
		state.statedModuleName = moduleName
		
		-- CECOMMPATCH - easypeasy display name check so everything isn't copypasted in modules.edb
		--  manual addition of display_name in modules.edb for those that need it (not all that many, turns out)
		local edb = EntityDB[state.name]
		if edb.displayname ~= nil then
			moduledisplayname = edb.displayname .. " (unbuilt)"
		else
			moduledisplayname = state.name .. " (unbuilt)"
		end
		
		--moduledisplayname = "Module Under Construction"
		state.statedBuildingName = buildingName
		state.moduleGhost = ghostModule

		if not state.complete then 
			send("rendGhostModelClassHandler",
					"odinRendererCreateMachineGhostModelRequest",
					SELF.id,
					state.moduleGhost,
					state.machinePos.x,
					state.machinePos.y,
					state.moduleWidth,						
					state.moduleHeight,
					state.moduleRotation,
					state.local_player,
					state.local_id,
					bh,
					SELF,
					false)
			
			send("rendGhostModelClassHandler",
					"odinRendererSetGhostModelName",
					SELF.id,
					moduledisplayname)
			
			if immediate then
				send(SELF, "ModuleConstructionComplete", nil)
			end
		end
	>>

	receive ModuleMoveCreationInformation( int modulePosX, int modulePosY, int moduleRotation )
	<<
		state.machinePos.x = modulePosX
		state.machinePos.y = modulePosY
		state.moduleRotation = moduleRotation
	>>
	receive SetParentBuilding ( gameObjectHandle pb )
	<<
		if VALUE_STORE["showModuleDebugConsole"] then printl("module", "setting parent building...") end
		state.parentBuilding = pb
	>>

	receive ClaimModule ( gameObjectHandle owner )
	<<
		local myTable = EntityDB[state.name]
		if myTable == nil then
			if VALUE_STORE["showModuleDebugConsole"] then printl("module", "module create: " .. state.name .. " not found...") end
 			return
		end
		local ownerName = query( owner, "getName" )[1]

		local tickerText = ownerName .. " has claimed a " .. state.name .. "."

		if myTable.icon and myTable.icon_skin then
		--	send("rendCommandManager", "odinRendererTickerMessage", tickerText, myTable.icon, myTable.icon_skin )
		else
		--	send("rendCommandManager", "odinRendererTickerMessage", tickerText, "carpentry_icon" , "ui\\orderIcons.xml" )
		end

		SELF.tags["unclaimed_module"] = nil
		state.owner = owner

	>>

	receive ModuleConstructionComplete( gameSimJobInstanceHandle ji )
	<<
		SELF.tags.under_construction = nil
		
		--Kinda hacky check to flip a game state variable
		if state.name == "Macroscope" then
			if not query("gameSession", "getSessionBool", "macroscopeBuilt")[1] then
				send("gameSession", "setSessionBool", "macroscopeBuilt", true)
			end
		end
		
		local myTable = EntityDB[state.name]
		if myTable == nil then
			if VALUE_STORE["showModuleDebugConsole"] then printl("module", "module create: " .. state.name .. " not found...") end
 			return
		end
		local tickerText = "Assembly of a " .. state.name .. " is complete."
		if myTable.icon and myTable.icon_skin then
		--	send("rendCommandManager", "odinRendererTickerMessage", tickerText, myTable.icon, myTable.icon_skin )
		else
		--	send("rendCommandManager", "odinRendererTickerMessage", tickerText, "carpentry_icon" , "ui\\orderIcons.xml" )
		end

		state.material_names = {}
		for i = 1, #state.materials do
			state.material_names[ i ] = query(state.materials[i],"getName")[1]
				
			send(state.materials[i], "DestroyedMessage")
			local resultROH = query( state.materials[i], "ROHQueryRequest" )
		
			if resultROH ~= nil then
				send("rendStaticPropClassHandler",
					"odinRendererDeleteStaticProp",
					resultROH[1])
			end
			
			destroyfromjob(state.materials[i], ji)
		end
		
		state.materials = {}
		state.repair_materials = {}
		local desired_height = 4
		if myTable.standalone == true then
			printl("buildings"," module setting desired height to 0");
			desired_height = 0
		end
		if SELF.tags["window"] == true or myTable.type == "buildingDecor" then
			printl("buildings"," module setting desired height to 0");
			desired_height = 0
		end
		if myTable.type == "interiorDecor" then
			desired_height = 4
		end
		
		send("rendGhostModelClassHandler", "odinRendererRemoveModuleTemplate", SELF.id);		-- pull this from the main scene graph.
		send("rendMachineClassHandler", "odinRendererCreateMachineMessage", 
			state.statedModuleName, 
			state.statedBuildingName,
			state.machinePos.x, 
			state.machinePos.y,
			state.moduleWidth,
			state.moduleHeight, 
			state.moduleRotation, 
			desired_height, 
			SELF,
			state.parentBuilding);
		send("rendMachineClassHandler", "odinRendererBindMachineToModuleTemplate", SELF.id);

		state.renderHandle = SELF.id
		if SELF.tags["broken"] == true then
			-- note: this is a pathing check
			print("brokentest", "Broken from ModuleConstructionComplete for " .. state.name )
			send("rendMachineClassHandler", "odinRendererMachineSetBrokenMessage", state.renderHandle, true);
		end
		state.complete = true
		
		if myTable.default_animation then
			send("rendMachineClassHandler", "odinRendererSetMachineDefaultAnimation", SELF.id, myTable.default_animation)
		end

 		send("rendInteractiveObjectClassHandler",
			"odinRendererSetAutomaticInteractionMessage",
			state.renderHandle,
			"menu")

		if state.parentBuilding ~= nil then
			send("gameWorkshopManager", "BuildingRegisterModule", state.parentBuilding, SELF, state.name, true);
			result = query("gameWorkshopManager", "getBuildingModulesGameSide", state.parentBuilding);
			local resultModuleTable = result[1]
			for k,v in pairs(resultModuleTable) do
				printl("buildings", "Module ID: " .. k .. " created with tags: " .. tostring( v.tags ) )
			end
			
			-- Recalculate the building quality now that the module is set up.
			send(state.parentBuilding, "recalculateQuality")
		else
				send("rendInteractiveObjectClassHandler",
					"odinRendererClearInteractions",
					state.renderHandle)
			
			send("rendInteractiveObjectClassHandler",
                    "odinRendererAddInteractions",
                    state.renderHandle,
                    "Dismantle Object",
                    "Dismantle Objects",
                    "Dismantle Objects",
                    "Dismantle Objects",
                    "hammer_icon",
                    "construction",
                    "Hammer Wood E",
					true,true)
		end
		
 	>>

	receive InteractiveMessage( string message )
	<<
		if message == "demolish" then
			send(SELF, "Demolition")
		end
	>>

	receive Demolition()
	<<
		printl("buildings","received Demolition order for: " .. state.statedModuleName  )
		
		if not state.complete then			
			send("rendGhostModelClassHandler", "odinRendererDeleteGhostModel", SELF.id, true)			
			send("rendGhostModelClassHandler", "odinRendererDeleteModuleTemplate", SELF.id)
			send("gameSpatialDictionary", "gridRemoveObject", SELF)
			destroy(SELF)
			return
		end
		
		send("gameWorkshopManager",
			"BuildingRegisterModule", state.parentBuilding, SELF, state.name, false)
		
		printl("buildings","creating dismantle job for " .. state.statedModuleName  )
		
		send("rendMachineClassHandler",
			 "odinRendererMachineExpressionMessage",
			 state.renderHandle,
			 "machine_thought",
			 "not_work")
		
		-- sfx for demolition:
		send("rendInteractiveObjectClassHandler",
			"odinRendererPlaySFXOnInteractive",
			SELF.id,
			"key_click")

		local assignmentResults = query("gameBlackboard",
								"gameObjectNewAssignmentMessage",
								SELF,
								"Demolish Module",
								"construction",
								"construction")
		
--		state.curDestructionAssignment = assignmentResults[1]

		send("gameBlackboard",
			"gameObjectNewJobToAssignment",
			assignmentResults[1],
			SELF,
			"Disassemble Module",
			"module",
			true )
		
		if state.parentBuilding ~= nil then
			send("gameWorkshopManager", "BuildingDeregisterModule", state.parentBuilding, SELF);
		end
	>>
	receive declareBroken()
	<<
		-- Note: this is not "broken" in terms of gameplay, this is broken in terms of pathfinding.
		SELF.tags.broken = true
		send("rendMachineClassHandler", "odinRendererMachineSetBrokenMessage", state.renderHandle, true);
	>>
	
	receive update()
	<<

	>>

	respond getModuleName ()
	<<
		return "nameResult", state.name
	>>

	receive ModuleItemTransform ( string name, string processingJobName, int processingJobID )
	<<
		state.processingName = name
		state.processingJobName = processingJobName
		state.processingTime = 100
		state.processingJobID = processingJobID
		state.processing = true
		SELF.tags["processing"] = true

		local jobInfo = EntityDB[processingJobName]
		if jobInfo ~= nil then
			local success = false
			if jobInfo.animations then
				for key,value in pairs(jobInfo.animations) do
					if key == state.name then
						state.processingTime = value.time
					end
				end			
			end
		end

		-- pull EDB, do stuff.
		local myTable = EntityDB[state.name]
		if myTable == nil then
			printl("buildings", " WARNING: module edb for " .. state.name .. " not found...");
			return
		end

		send("rendMachineClassHandler",
			 "odinRendererMachineExpressionMessage",
			 state.renderHandle,
			 "machine_thought",
			 "work")
	>>

	respond getWorkshopProduct()
	<<
	>>

	receive CompleteJob( int index )
	<<
		if VALUE_STORE["showModuleDebugConsole"] then printl("module", "module type: " .. state.name .. " called CompleteJob" ) end
 		SELF.tags["processing"] = false
		send(state.parentBuilding, "CompleteJob", index )
	>>

	receive CompleteJobByName ( string name )
	<<
		if VALUE_STORE["showModuleDebugConsole"] then printl("module", "module type: " .. state.name .. " called CompleteJobByName") end
 		SELF.tags["processing"] = false
		send(state.parentBuilding, "CompleteJob", state.processingJobID )				-- hack
	>>

	respond getModuleProcessingJob()
	<<
		return "nameResult", state.processingJobName
	>>

	respond getModuleProcessingName()
	<<
		return "nameResult", state.processingName
	>>
	
	receive ModuleUsed()
	<<
		-- maybe replacing this with upkeep system.
		-- yes, yes we are.
		
		--[[
		state.timesUsed = state.timesUsed + 1

		if state.timesUsed >= state.usesUntilDamaged and not SELF.tags["damaged"] then
			
			-- Can't use this for repair materials anymore due to boxed module system as of 45A. -dgb
			-- Will now set repair materials in modules.edb entries for modules, or default to something reasonable.
			--local result = query("gameUtilitySingleton", "odinGetRandomModuleRepairPart", state.name)
						
			send(SELF,"spawnDebris")
			
			local tag_to_entity = {
				timber = "logs",
				rough_stone_block = "rough_stone_block",
				planks = "planks",
				masonry = "bricks",
				bricks = "bricks",
				bolt_of_cloth = "bolt_of_cloth",
				ingot_of_iron = "iron_ingots",
				ingot_of_copper = "copper_ingots",
				ingot_of_brass = "brass_ingots",
				ingot_of_zinc = "zinc_ingots",
				ingot_of_gold = "gold_ingots",
				brass_cogs = "brass_cogs",
				iron_pipes = "iron_pipes",
				copper_pipes = "copper_pipes",
				iron_plates = "iron_plates",
				copper_plates = "copper_plates",
				glass_panes = "glass_panes",
			}
			
			-- NOTE: repair material is a tag. This should be in an edb somewhere Responsible.
			local entityInfo = EntityDB[state.name]
			local materialTag = ""
			if entityInfo.repair_materials then
				materialTag = entityInfo.repair_materials[ rand(1,#entityInfo.repair_materials) ]
			elseif entityInfo.tier then
				-- go by tier
				local repair_tier = rand(1,entityInfo.tier)
				if repair_tier > 4 then repair_tier = 4 end
				local repair_materials = {
					[1] = { "timber", "rough_stone_block" },
					[2] = { "planks", "masonry", "planks" },
					[3] = { "iron_ingots", "copper_ingots", "iron_plates", "iron_pipes"},
					[4] = { "brass_cogs", "copper_pipes", "copper_plates", "glass_panes" },
				}
				local choices = repair_materials[ rand(1,repair_tier) ]
				materialTag = choices[ rand(1,#choices) ]
			else
				-- uh, planks I guess?
				materialTag = "planks"
			end
			
			local materialEntityName = tag_to_entity[materialTag]
			
			if not materialEntityName then
				-- problems!
				printl("module", "PROBLEMS! no materialEntityName")
				return	
			end
			
			local materialName = EntityDB[materialEntityName].display_name
			
			if not materialName then
				-- problems!
				printl("module", "PROBLEMS! no materialName")
				return
			end
			
			local s = "A " .. state.name .. " has become damaged due to overuse and will require repair (requires: " .. materialName .. ")."
			send("rendCommandManager", "odinRendererTickerMessage", s , "broken_cog", "ui\\thoughtIcons.xml" )
			
			
			SELF.tags["damaged"] = true
			SELF.tags["inoperative"] = true -- unifying the damage states for simplicity
			
			local damageSound = "Grenade Blast"
			local damagePFX = "MidSplosion"
			
			if SELF.tags["bed"] then
				damageSound = "Break (Dredmor)"
				damagePFX = "TreeChoppedDown"
			end

			send("rendInteractiveObjectClassHandler",
					"odinRendererPlaySFXOnInteractive",
					state.renderHandle,
					damageSound )
			
			send("rendCommandManager",
					"odinRendererCreateParticleSystemMessage",
					damagePFX,
					state.machinePos.x,
					state.machinePos.y )
			 
			-- You broke it! Make a repair job.
			send("rendCommandManager",
				"odinRendererCreateParticleSystemMessage",
				"Small Beacon",
				state.machinePos.x,
				state.machinePos.y)
			
			-- TODO give this to workshop workcrew, if applicable?
			-- though, note, workshop will often shuffle contents of its assignment
			-- thus possibly pushing off any outside added job
			
			local assignment = query("gameBlackboard",
						"gameObjectNewAssignmentMessage",
						SELF,
						"Repair Module",
						"construction",
						"construction")[1] 

			send( "gameBlackboard", "gameObjectNewJobToAssignmentWithJobTag",
					assignment,
					SELF,
					"Gather Module Repair Materials",
					"module",
					true,
					materialTag)
               
			send( "gameBlackboard", "gameObjectNewJobToAssignment",
					assignment,
					SELF,
					"Repair Module",
					"module",
					true )
			
			send("rendMachineClassHandler",
					"odinRendererMachineExpressionMessage",
					state.renderHandle,
					"machine_thought",
					"broken_cog")
			 
		end]]
	>>
	
	receive ModuleRepairComplete(gameSimJobInstanceHandle ji)
	<<
		SELF.tags["damaged"] = nil
		SELF.tags["inoperative"] = nil
		state.timesUsed = 0
		
		send("rendCommandManager",
			"odinRendererCreateParticleSystemMessage",
			"Sparkle",
			state.machinePos.x,
			state.machinePos.y)
			
		send("rendMachineClassHandler",
				"odinRendererMachineExpressionMessage",
				state.renderHandle,
				"",
				"")

		state.repair_material_dropped_off = false

		-- destroy repair material!

		for i = 1, #state.repair_materials do
			send(state.repair_materials[i], "DestroyedMessage")
			local resultROH = query( state.repair_materials[i], "ROHQueryRequest" )
		
			if resultROH ~= nil then
				send("rendStaticPropClassHandler",
					"odinRendererDeleteStaticProp",
					resultROH[1])
			end
			
			destroyfromjob(state.repair_materials[i], ji)
		end
	>>
	
	receive ModuleSetCreator( gameObjectHandle creator, string name )
	<<
		-- was a dangling hook for this, so testing if it can do what was intended.
		---printl("DAVID", state.statedModuleName .. " received ModuleSetCreator from " .. name)
		-- it works! Nice.
		
		-- someone do something cool here.
	>>

	receive JobCancelledMessage ( gameSimJobInstanceHandle ji )
	<<
		if ji.name == "Finish Production" then
			state.processing = false
			SELF.tags["processing"] = nil
			send("rendMachineClassHandler",
				"odinRendererMachineExpressionMessage",
				state.renderHandle,
				"",
				"")
			
			send("rendMachineClassHandler", "odinRendererResetMachineAnimation", SELF.id);
		elseif ji.name == "Construct Module" or ji.name == "Construct Module (in building)" then
				send("rendGhostModelClassHandler", "odinRendererDeleteGhostModel", SELF.id, false);
				
		elseif ji.name == "Disassemble Module" then
			printl("buildings", tostring(SELF.id) .. " got JobCancelledMessage for Disassemble Module")
			-- re-register self
			if state.parentBuilding ~= nil then
				send("gameWorkshopManager", "BuildingRegisterModule", state.parentBuilding, SELF, state.name, true);
			end
		end
	>>
	
	respond getParentBuilding()
	<<
		return "building", state.parentBuilding
	>>
	
	receive spawnDebris()
	<<
		local debrisAmount = 3
		local debrisTable = { --"Machine Rubbish",
						 "Paper Rubbish",
						 "Random Rubbish",
						 --"Brick Debris, Small",
						 --"Broken Glass",
						 "Wooden Debris, Small" }
		
		--Machine Rubbish
		--
		-- TODO: pull above values from entityDB
			
		for i=1, debrisAmount do
				
			local handle = query("scriptManager",
				"scriptCreateGameObjectRequest",
				"clearable",
				{ legacyString = debrisTable[ rand(1,#debrisTable) ]} )[1]
			
			local px = state.position.x + rand(-2,2)
			local py = state.position.y + rand(-2,2) 
			
			send( handle,
				"GameObjectPlace",
				px,
				py)
		
			local damagePFX = "DustPuffLarge"
			send("rendCommandManager",
					"odinRendererCreateParticleSystemMessage",
					damagePFX,
					px,
					py )
			-- TODO: do dustcloud on spawnpos.
		end
	>>
	
	receive addTag( string name )
	<<
		SELF.tags[name] = true
	>>
	
     respond getTags()
     <<
          return "getTagsResponse", SELF.tags
     >>
	
	receive removeTag( string name )
	<<
		SELF.tags[name] = nil
	>>
	
	respond getRenderhandle()
	<<
		return "getRenderhandleResponse", state.renderHandle
	>>
>>