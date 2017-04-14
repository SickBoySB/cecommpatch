gameobject "clearable" inherit "object_damage"
<<
	local 
	<<
          function resetClearableDescription()
               local description = "Just a thing on the ground."
               entityData = EntitiesByType["clearable"][state.clearableName]
               if entityData.description then
                    description = entityData.description
               end
               
               if SELF.tags["devCoords"] then
                    local civ = query( "gameSpatialDictionary", "gridGetCivilization", state.position )[1]
                    description = description .. " / position x= " .. tostring(state.position.x) .. ", y= " .. tostring(state.position.y) .. " / civ= " .. tostring(civ)
               end
			
			--description = description .. " Remove by using the 'Clear Terrain' command."
               
               send("rendInteractiveObjectClassHandler",
                    "odinRendererBindTooltip",
                    state.renderHandle,
                    "ui//tooltips//groundItemTooltipDetailed.xml",
                    entityData.name,
                    description )
          end
          
          function resetClearableModel()
               if state.modelPlaced then
                    send("rendStaticPropClassHandler", "odinRendererDeleteStaticProp", SELF.id)
               end
               state.modelPlaced = true
               
               local entityData = EntitiesByType["clearable"][state.clearableName]
               
               if entityData.model ~= nil then
                    -- for single model given
                    send("rendStaticPropClassHandler", "odinRendererCreateStaticPropRequest", SELF,
                         entityData.model, state.position.x, state.position.y)
               else
                    -- random index from models list
                    r =  rand(1, #entityData.models) 
                    send("rendStaticPropClassHandler", "odinRendererCreateStaticPropRequest", SELF,
                         entityData.models[ r ], state.position.x, state.position.y)
               end 
               send("rendStaticPropClassHandler", "odinRendererRotateStaticProp", SELF.id, state.modelRotation, 0.25)
               resetClearableDescription()
          end
          
          function resetClearableCommands()
               send("rendInteractiveObjectClassHandler", "odinRendererClearInteractions", state.renderHandle)
               if not state.addedJob then
                    send("rendInteractiveObjectClassHandler", "odinRendererAddInteractions", state.renderHandle,
                         "Clear " .. EntitiesByType["clearable"][state.clearableName].name,
                         "Clear Terrain",
                              "Clear Terrain",
                              "Clear Terrain",
                              "clearable",
							"terrain",
							"Dirt",
							true,false)
			end
          end

	>>

	state
	<<
		gameGridPosition position
		int renderHandle
		string clearableName
		int modelRotation
		bool addedJob
          bool modelPlaced
		gameSimAssignmentHandle assignment
	>>

	receive Create(stringstringMapHandle init)
	<<
		local objectName = init["legacyString"]
		state.assignment = nil

		addedJob = false
          
		state.position.x = -1
		state.position.y = -1
		state.clearableName = objectName
          
		entityData = EntityDB[state.clearableName]
		if entityData == nil then
			ScriptError("Clearable", "Clearable place Failed: Couldn't find " .. objectName)
		end
          
		SELF.tags={}
		if entityData.tags then
			for k,v in pairs(entityData.tags) do
				SELF.tags[v] = true	
			end
			SELF.tags.clearable = true
		else
			SELF.tags = { "clearable" }
		end
          
          if SELF.tags["no_rotation"] then
               state.modelRotation = 0
          else
               state.modelRotation = rand(0, 359)
          end
          
          state.modelPlaced = false
		
		if init.timer then state.rotTimer = tonumber(init.timer) end
          
		ready()
		sleep()
	>>

	receive GameObjectPlace( int x, int y ) 
	<<
		state.position.x = x
		state.position.y = y
		state.renderHandle = SELF.id

          resetClearableModel()
          resetClearableCommands()
		
		send( "gameSpatialDictionary", "registerSpatialMapString", SELF, "c", "c", true )		 
		send("gameSpatialDictionary", "gridAddObjectTo", SELF, state.position)
		
          local entityData = EntitiesByType["clearable"][state.clearableName]
          
		if entityData[obstruction] == true then
			send("gameSpatialDictionary", "gridSetPassable", state.position, "landscape", true)
		else
               -- nope
			--send("gameSpatialDictionary", "gridSetPassable", state.position, "landscape", false)
		end
          
          -- HAX add particles
          if entityData.pfx then
               send("rendCommandManager", "odinRendererCreateParticleSystemMessage",
                    entityData.pfx,
                    state.position.x,
                    state.position.y)
          end
          
          if SELF.tags["rots"] then
			if state.rotTimer then
				-- done.
			elseif entityData.decayTimeDays then
				state.rotTimer = entityData.decayTimeDays * EntityDB["WorldStats"]["dayNightCycleTenthSeconds"]
			elseif entityData.decayTimeSeconds then 
				state.rotTimer = entityData.decayTimeSeconds * 10 -- for tenths of second
			end
               wake()
			
		elseif SELF.tags.grenade and SELF.tags.primed then
			if state.rotTimer then
				-- done.
			else
				-- default
				state.rotTimer = 17
			end
			wake()
          end

		--local description = "position x= " .. tostring(state.position.x) .. ", y= " .. tostring(state.position.y)
		--printl("dev coords", description);
		
		if SELF.tags["devCoords"] then
			send("gameSession", "setSessionBool", "airdropOverride", true)
			send("gameSession", "setSessionInt", "airdropX", x)
			send("gameSession", "setSessionInt", "airdropY", y)
		end
		
		local iswater = query( "gameSpatialDictionary",
							"gridHasSpatialTag",
							state.position,
							"water" )[1]
		
		if iswater then SELF.tags.placed_in_water = true end
	>>

	receive GameObjectAddInstance( int x, int y ) 
	<<
          entityData = EntitiesByType["clearable"][state.clearableName]
		send("rendStaticPropClassHandler", "odinRendererCreateStaticPropRequest", SELF,entityData.model, x, y)
	>>

	receive gameFogOfWarExplored (int x, int y)
     <<
		if SELF.tags.bandit_campfire then
			
			local eventQ = query("gameSimEventManager", "startEvent", "bandit_camp_warning", {}, {})
			send(eventQ,"registerSubject", SELF )
			send(eventQ,"registerPosition", state.position )
		end
	>>

	receive Update()
	<<
		if state.rotTimer then
			if not SELF.tags["occupied"] then
				state.rotTimer = state.rotTimer - 1
				if state.rotTimer <= 0 then
					if SELF.tags.grenade then
						
						local handle = query("scriptManager",
									"scriptCreateGameObjectRequest",
									"explosion",
									{legacyString= "Medium Explosion" })[1]
				
						send(handle,
							"GameObjectPlace",
							state.position.x,
							state.position.y)
				
						send(SELF, "Clear", nil)
					else
						send(SELF, "Clear", nil)
					end
					
				elseif state.rotTimer % 3 == 0 and SELF.tags.primed then
					send("rendCommandManager",
						"odinRendererCreateParticleSystemMessage",
						"DustPuffV1",
						state.position.x,
						state.position.y)
				end
			end
		end
	>>

	respond gridGetPosition()
	<<
		return "reportedPosition", state.position
	>>

	respond isObstruction()
	<<
		return "obstructionResult", false
	>>

	respond gridReportPosition()
	<<
		return "gridReportedPosition", state.position
	>>

	receive AssignmentCancelledMessage( gameSimAssignmentHandle assignment )
	<<
		state.addedJob = false
	>>

	receive InteractiveMessage( string messagereceived )
	<<
		send(SELF,"handleInteractiveMessage", messagereceived, nil)
     >>
     
     receive InteractiveMessageWithAssignment( string messagereceived, gameSimAssignmentHandle assignment )
	<<
          send(SELF,"handleInteractiveMessage", messagereceived, assignment)
     >>
	
	receive handleInteractiveMessage( string messagereceived, gameSimAssignmentHandle assignment )
	<<
		if not state.addedJob then
               if messagereceived == "Clear Terrain" then
                    if not assignment then
					
					assignment = query("gameBlackboard",
								"gameObjectNewAssignmentMessage",
								SELF,
								"Clear Terrain",
								"construction",
								"construction")[1]
					
				end
				
                    state.assignment = assignment
                    state.addedJob = true
				
				send("rendBeaconClassHandler",
					"CreateAssignmentBeacon",
					state.assignment,
					"jobshovel",
					"ui\\thoughtIcons.xml",
					"ui\\thoughtIconsGray.xml",
					3)
				
				send("rendBeaconClassHandler",
					"AddEntityToAssignmentBeacon",
					state.assignment,
					SELF.id)
				
                    send( "gameBlackboard",
                         "gameObjectNewJobToAssignment",
                         state.assignment,
                         SELF,
                         "Clear Terrain",
                         "clearable",
                         true )
				
				resetClearableCommands()
				
               elseif messagereceived == "Chop Down" then
				
				if not assignment then
					assignment = query("gameBlackboard",
									"gameObjectNewAssignmentMessage",
									SELF,
									"Chop Tree",
									"chopping",
									"chopping")[1]
				end
				
				state.assignment = assignment
                    state.addedJob = true
				
				send("rendBeaconClassHandler",
					"CreateAssignmentBeacon",
					state.assignment,
					"jobaxe",
					"ui\\thoughtIcons.xml",
					"ui\\thoughtIconsGray.xml",
					3)
				
				send("rendBeaconClassHandler",
					"AddEntityToAssignmentBeacon",
					state.assignment,
					SELF.id)
				
				send("gameBlackboard",
					"gameObjectNewJobToAssignment",
					state.assignment,
					SELF,
					"Chop Tree",
					"tree",
					true)
				
               elseif messagereceived == "Clear Area" then
                    results = query("gameSpatialDictionary", "allObjectsInRadiusRequest", state.position, 4, true)
                    if results ~= nil then
                        send(results[1], "InteractiveMessage", "Clear Clearables Only")
                    end
               elseif messagereceived == "Clear Clearables Only" then
                    send(SELF,"InteractiveMessage", "Clear Terrain")
               end
          end
	>>
	
	receive ClearNature( gameSimAssignmentHandle assignment ) 
	<<
		-- Generic clear command for all nature objects.
		if not state.addedJob then
			send(SELF, "handleInteractiveMessage", "Clear Terrain", assignment)
		end
	>>

	receive Harvest( gameObjectHandle harvester, gameSimJobInstanceHandle ji )
	<<
		if state.nestOwner ~= nil then
			send(state.nestOwner,"NestCleared")
		end
		
		if query("gameSession", "getSessionBool", "trackClearables")[1] == true then
			local director = query("gameSession","getSessiongOH","event_director_eldritch")[1]
			send(director,"incKeyInt", "cleared", 1)
		end
	
		-- rolling receive Clear into this
		-- Put some music in the dynamics.
		incMusic(3,10);
		-- hey maybe spawn some junk
		
		local data = EntityDB[state.clearableName]
		if data.commodityAmount and data.commodityOutput then
			for i=1, data.commodityAmount do
				local handle  = query( "scriptManager",
								  "scriptCreateGameObjectRequest",
								  "item",
								  { legacyString = data.commodityOutput })[1]
				
				local civ = query("gameSpatialDictionary","gridGetCivilization",state.position)[1]
				--if civ == 0 then
					send(handle,"ClaimItem")
				--else
				--	send(handle,"ForbidItem")
				--end		
				
				send(handle,
					"GameObjectPlace",
					state.position.x,
					state.position.y)
			end
			
		elseif data.spawnAfterClear then
			local spawntable = data.spawnAfterClear
			for k, v in pairs(spawntable) do
				local handle  = query( "scriptManager",
								  "scriptCreateGameObjectRequest",
								  v.entityType,
								  { legacyString = v.entityName })[1]
				
				if v.entityType == "item" then
					local civ = query("gameSpatialDictionary","gridGetCivilization",state.position)[1]
					if civ == 0 then
						send(handle,"ClaimItem")
					else
						send(handle,"ForbidItem")
					end		
				end
				
				send(handle,
					"GameObjectPlace",
					state.position.x,
					state.position.y)
			end
		end
		
		if ji then
			if SELF.tags["rock"] or SELF.tags["grass"] then
				-- if a human is doing this job and not an explosion
				local worldstats = EntityDB["WorldStats"]
				if query("gameSession","getSessionBool","creepyCultEventSeen")[1] then
					
					if rand(1, worldstats["artifactFindChance"] ) == 1 then
						-- hey, what's under this rock/plant? it's a strange artifact!
						
						local eventQ = query("gameSimEventManager",
									"startEvent",
									"spawn_artifact",
									{},
									{ reason = "clearing terrain" } )[1] 
							
						send(eventQ,"registerSubject",harvester)
						send(eventQ,"registerPosition", state.position)
					end
					
				else
					if rand(1,300 ) == 1 then
						-- hey, what's under this rock/plant? it's a strange artifact!
						local eventQ = query("gameSimEventManager",
									"startEvent",
									"spawn_artifact",
									{},
									{ reason = "clearing terrain" } )[1] 
							
						send(eventQ,"registerSubject",harvester)
						send(eventQ,"registerPosition", state.position)
					end
				end
			end
			
		end
		
		if SELF.tags.grass then
			send("gameSession","incSessionInt","grassClearedCount",1)
			local grasscount = query("gameSession","getSessionInt","grassClearedCount")[1]
			if grasscount >= 50 and
				not query("gameSession","getSessionBool","clearedLotsOfGrass")[1] then
				
				send("gameSession","setSessionBool","clearedLotsOfGrass", true)
			end
			send("gameSession", "incSteamStat", "stat_grass_cleared_count", 1)
		end
		
		-- Finaeely, destroy myself.
          send("gameBlackboard", "gameObjectRemoveTargetingJobs", SELF, ji)
		send("rendStaticPropClassHandler", "odinRendererDeleteStaticProp", SELF.id)
		send("gameSpatialDictionary", "gridRemoveObject", SELF);
		destroyfromjob(SELF, ji )
	>>
	
	receive ClearViolently( gameSimJobInstanceHandle ji, gameObjectHandle damagingObject )
	<<
		send(SELF,"Clear",ji)
	>>
	
	receive Clear(gameSimJobInstanceHandle ji )
	<<
		-- Put some music in the dynamics.
		incMusic(3,10)
		-- hey maybe spawn some junk
		if EntityDB[state.clearableName].spawnAfterClear then
			local spawntable = EntityDB[state.clearableName].spawnAfterClear
			for k, v in pairs(spawntable) do
				local handle  = query( "scriptManager","scriptCreateGameObjectRequest", v.entityType, { legacyString = v.entityName })[1]
				send(handle,
					"GameObjectPlace",
					state.position.x,
					state.position.y)
			end
		end
		
		-- Finaeely, destroy myself.
          send("gameBlackboard", "gameObjectRemoveTargetingJobs", SELF, ji)
		send("rendStaticPropClassHandler", "odinRendererDeleteStaticProp", SELF.id)
		send("gameSpatialDictionary", "gridRemoveObject", SELF);
		destroyfromjob(SELF, ji )
	>>
     
     respond getAnimLoopsToClear()
     <<
          entityData = EntitiesByType["clearable"][state.clearableName]
          if entityData.animLoopsToClear then
               return "animLoopsToClear", entityData.animLoopsToClear
          else
               return "animLoopsToClear", 1
          end
     >>
     
     respond getRenderHandle()
     <<
          if state.renderHandle then
               return "getRenderHandle", state.renderHandle
          else
               return "getRenderHandle", nil
          end
     >>
	
	receive JobCancelledMessage( gameSimJobInstanceHandle ji )
	<<
          state.addedJob = false
          state.assignment = nil
          resetClearableModel()
          resetClearableCommands()
	>>
	
	respond ROHQueryRequest()
	<<
		return "ROHQueryReply", state.renderHandle 
	>>
	
	respond getChopAnimation()
	<<
		if SELF.tags.logs then
			return "getChopAnimationMessage", "mine_node"
		end
		return "getChopAnimationMessage", nil
	>>
	
	respond getHarvestCommodityOutput()
	<<
		return "getHarvestCommodityOutputMessage", EntityDB[ state.entityName ].commodityOutput
	>>
	
	respond getHarvestCommodityAmount()
	<<
		return "getHarvestCommodityAmountMessage", EntityDB[ state.entityName].commodityAmount
	>>
     
     respond getName()
	<<
		return "clearableNameRespond", state.clearableName
	>>
>>
