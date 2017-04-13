gameobject "outsider" inherit "ai_agent"
<<
	local
	<<
          function outsider_doOneSecondUpdate()
               state.AI.ints["emoteTimer"] = state.AI.ints["emoteTimer"] + 1
               
               -- normalize
               if state.AI.ints["morale"] < 0 then
                    state.AI.ints["morale"] = 0
               elseif state.AI.ints["morale"] > 100 then
                    state.AI.ints["morale"] = 100
               else
                    state.AI.ints["morale"] = state.AI.ints["morale"] +1
               end
               
               if SELF.tags["fleeing"] and state.AI.ints["morale"] >= 50 then 
                    SELF.tags["fleeing"] = false
               end               
          end
	>>

	state
	<<
		gameObjectHandle group
		gameAIAttributes AI
		gameGridPosition spawnLocation
		gameGridPosition wanderDestination
		table entityData
		table traits
		string entityName
		string animSet
		string nation
		int renderHandle
		bool asleep
	>>
	
	receive Create( stringstringMapHandle init )
	<<
		local entityName = init["legacyString"]
          if not entityName then
               return
          end
          state.entityName = entityName
		
		local entityData = EntityDB[ state.entityName ]
		if not entityData then
			printl("ai_agent", "outsider type not found")
			return
		end
          state.entityData = entityData
	    
          state.AI.strs["mood"] = "content"
		
          if rand(0,100) < 50 then 
			state.AI.strs["gender"] = "male"
               state.AI.strs["firstName"] = maleFirstNames[ rand(1, #maleFirstNames ) ]
		else
			state.AI.strs["gender"] = "female"
               state.AI.strs["firstName"] =  femaleFirstNames[ rand(1, #femaleFirstNames ) ]
		end
		
		local lastname = ""
		local lastA = lastnamePrefixes[ rand(1, #lastnamePrefixes ) ]
		local lastB = lastnameRoots[ rand(1, #lastnameRoots ) ]
		local lastC = lastnameSuffixes[ rand(1, #lastnameSuffixes ) ]

		local temprand = rand(1,6)
		if( temprand == 1 ) then
			lastname = lastA .. lastB
		elseif (temprand == 2) then
			lastname = lastB .. lastC
		elseif (temprand == 3) then
			lastname = lastB
		elseif (temprand == 4) then
			--state.AI.strs["lastName"] = lastnames[ rand(1, #lastnames ) ]
			lastname = lastA .. lastC
		else
			--state.AI.strs["lastName"] = lastA .. lastB .. lastC
			lastname = lastnames[ rand(1, #lastnames ) ]
		end
		
		local inspector_titles = {
			"Inspector", "Sergeant", "Commissionair", "Constable", "Deputy Chief Constable",
			"Chief Constable", "Assistant Deputy Constable", "Deputy Assistant Constable",
			"Deputy Inspector", "Assistant Inspector", "Deputy Assistant Inspector", "Assistant Deputy Inspector",
			"Deputy Sergeant", "Assisant Sergeant", "Deputy Assistant Sergeant", "Assistant Deputy Sergeant",
			"Deputy Purger", "Assistant Purger", "Sergeant-Purger", "Inspector-Purger",
		}
		
		local inspector_titles_high_rank = {
			"Commander", "Assistant Commander", "Deputy Commander", "Purger", "Chief Purger",
			"Deputy Assistant Commissioner", "Assistant Commissioner", "Deputy Commissioner",
			"Deputy Assistant to the Master Purifier",
			"Assistant Secretary to the Master Purifier",
			"Assistant High All-Purger", "Deputy Witch-hunter General", "Assistant Witch-hunter General",
		}
		
		if state.entityName == "Occult Inspector" then
			state.AI.strs["firstName"] = inspector_titles[ rand(1,#inspector_titles) ]
		elseif state.entityName == "Occult Sergeant" then
			state.AI.strs["firstName"] = inspector_titles_high_rank[ rand(1,#inspector_titles_high_rank) ]
		end
		
		state.AI.strs["lastName"] = string.upper( lastname:sub(1,1) ) .. lastname:sub(2, #lastname)
		state.AI.name = state.AI.strs.firstName .. " " .. state.AI.strs.lastName
		
		
		local isArmoured = false
		local variants = {[1]="A", [2]="B", [3]="C", [4]="D"}
		state.AI.strs["variant"] = variants[rand(1,4)]
		
		state.AI.strs["citizenClass"] = entityName
		local models = getModelsForClass( state.AI.strs["citizenClass"],
									state.AI.strs["gender"],
									state.AI.strs["variant"] )

		state.animSet = models["animationSet"]


		local hair = models.hairModel
		local hat = entityData.hatModel
		if hat and hat ~= "" then
			-- don't combine hat + hair.
			hair = ""
		end
		
          send("rendOdinCharacterClassHandler",
			"odinRendererCreateCitizen", 
			SELF, 
			models["torsoModel"], 
			models["headModel"],
			hair, --models["hairModel"], -- "models/hats/phrygiancap.upm", --models["hairModel"], 
			hat, --models["hatModel"], 
			models["animationSet"], 0, 0 )
		
		state.headModel = models["headModel"]
		state.hairModel = models["hairModel"]
		state.hatModel = entityData.hatModel
		models.hatModel = state.hatModel
		
		state.models = models

		send("rendOdinCharacterClassHandler",
			"odinRendererFaceCharacter", 
				state.renderHandle, 
				state.AI.position.orientationX,
				state.AI.position.orientationY )

		local humanstats = EntityDB["HumanStats"]
		local worldstats = EntityDB["WorldStats"]

		SELF.tags = {}
          for i,tag in ipairs(entityData.tags) do
			SELF.tags[tag] = true	
		end
		
		SELF.tags["conversible"] = nil

		if isArmoured then
			SELF.tags["armoured"] = true
		end
		SELF.tags.hostile_vs_bandits = true
		SELF.tags.hostile_vs_fishpeople = true
		
		state.AI.ints["grenades"] = state.AI.ints["grenadesMax"] -- and grenades, why not.

          -- START ai_damage required stats
		state.AI.ints["healthMax"] = humanstats["healthMax"]
          state.AI.ints["healthTimer"] = 3 -- in seconds, per 1 point 
          state.AI.ints["fire_timer"] = 10
          state.AI.ints["health"] = state.AI.ints["healthMax"]
		state.AI.ints["numAfflictions"] = 0
		
          -- END ai_damage required stats
          
          state.AI.ints["wall_attacks"] = 0
           
          state.AI.ints["morale"] = 100
          state.AI.ints["corpse_timer"] = humanstats.corpseRotTimeDays * worldstats["dayNightCycleSeconds"] * 10 -- gameticks
		state.AI.ints.corpse_vermin_spawn_time_start = div(state.AI.ints.corpse_timer,2)
          state.AI.walkTicks = 3
          state.AI.ints["subGridWalkTicks"] = state.AI.walkTicks
          setposition(0,0)
		
		state.AI.ints["sightRadius"] = 16
		
          state.AI.walkTicks = 3 -- from citizen.go
		state.AI.ints["subGridWalkTicks"] = state.AI.walkTicks

          send( "rendOdinCharacterClassHandler", "odinRendererSetCharacterWalkTicks",  SELF.id, state.AI.walkTicks )

		state.AI.strs["occupancyMap"] = 
		".-.\\".. 
		"-C-\\"..
		".-.\\"
		
		state.AI.strs["occupancyMapRotate45"] =  
		".-.\\".. 
		"-C-\\"..
		".-.\\"
		
		send( "gameSpatialDictionary",
			"registerSpatialMapString",
			SELF,
			state.AI.strs["occupancyMap"],
			state.AI.strs["occupancyMapRotate45"],
			true )
    
          if state.AI.strs["gender"] == "female" then
			state.vocalID = "High00"
		else
			if( rand(0,100) < 50 ) then 
				state.vocalID = "Low00"
			else
				state.vocalID = "Mid00"
			end
		end

		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle,
			"name",
			state.AI.name)
		
		if SELF.tags.occult_inspector then
			state.AI.ints["healthMax"] = EntityDB[state.entityName].healthMax
			state.AI.ints["health"] = state.AI.ints["healthMax"]
			
			send("rendOdinCharacterClassHandler",
				"odinRendererSetCharacterCustomTooltipMessage",
				SELF.id,
				"ui\\tooltips\\occultInspectorTooltip.xml")
		else
			send("rendOdinCharacterClassHandler",
				"odinRendererSetCharacterCustomTooltipMessage",
				SELF.id,
				"ui\\tooltips\\clockworkianFriendlyTooltip.xml")
		end
          
		-- set up bits
		
		send("gameSpatialDictionary", "gameObjectClearBitfield", SELF)
		send("gameSpatialDictionary", "gameObjectClearHostileBit", SELF)
		
		-- I am a subject of the empire and player 1, so bits 0 and 4
		send("gameSpatialDictionary", "gameObjectAddBit", SELF, 0)
		send("gameSpatialDictionary", "gameObjectAddBit", SELF, 4)

		-- I am hostile to all other players.
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 1)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 2)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 3)
		
		
		send(SELF,"makeFriendly")
		--[[local myNation = query("gameSession","getSessiongOH",state.nation)[1]
		 
		if query(myNation, "isHostile")[1] then
			send(SELF,"makeHostile")
		elseif query(myNation, "isFriendly")[1] then
			
		else
			send(SELF,"makeNeutral")
		end]]
		
          ready()
	>>
	
	receive gameFogOfWarExplored(int x, int y )
	<<
		state.asleep = false
	>>

	receive GameObjectPlace( int x, int y ) 
	<<
		if not state.AI.bools.first_placement then 
          	
			local entity_data = EntityDB[ state.entityName ]
			
			send(SELF,"setWeapon","melee",entity_data.melee_weapons[ rand(1,#entity_data.melee_weapons) ] )
			send(SELF,"setWeapon","ranged",entity_data.ranged_weapons[ rand(1,#entity_data.ranged_weapons) ] )
			  
			state.homeLocation = gameGridPosition:new()
			state.homeLocation.x = x
			state.homeLocation.y = y
			
			state.AI.bools.first_placement = true
		end
	>>

	receive makeHostile()
	<<
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle,
			"name",
			state.AI.name)
		
		if SELF.tags.occult_inspector then
			send("rendOdinCharacterClassHandler",
				"odinRendererSetCharacterCustomTooltipMessage",
				SELF.id,
				"ui\\tooltips\\occultInspectorTooltip.xml")
		else
			send("rendOdinCharacterClassHandler",
				"odinRendererSetCharacterCustomTooltipMessage",
				SELF.id,
				"ui\\tooltips\\clockworkianFriendlyTooltip.xml")
		end
	>>
	
	receive makeFriendly()
	<<
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle,
			"name",
			state.AI.name)
		
		if SELF.tags.occult_inspector then
			send("rendOdinCharacterClassHandler",
				"odinRendererSetCharacterCustomTooltipMessage",
				SELF.id,
				"ui\\tooltips\\occultInspectorTooltip.xml")
		else
			send("rendOdinCharacterClassHandler",
				"odinRendererSetCharacterCustomTooltipMessage",
				SELF.id,
				"ui\\tooltips\\clockworkianFriendlyTooltip.xml")
		end
	>>
	
	receive makeNeutral()
	<<
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle,
			"name",
			state.AI.name)
		
		if SELF.tags.occult_inspector then
			send("rendOdinCharacterClassHandler",
				"odinRendererSetCharacterCustomTooltipMessage",
				SELF.id,
				"ui\\tooltips\\occultInspectorTooltip.xml")
		else
			send("rendOdinCharacterClassHandler",
				"odinRendererSetCharacterCustomTooltipMessage",
				SELF.id,
				"ui\\tooltips\\clockworkianFriendlyTooltip.xml")
		end
	>>
	
	receive SleepMessage()
	<<
		state.asleep = true
	>>

     receive Update()
     <<
		tooltip_refresh_from_save()
		
  		if state.AI.thinkLocked then
               return
          end
		
          if state.AI.ints.updateTimer % 10 == 0 then
               outsider_doOneSecondUpdate()
			if (SELF.tags.exit_map or SELF.tags.exit_map_run) and SELF.tags.occult_inspector then
				local citizenNearbyResults = query("gameSpatialDictionary","isObjectInRadiusWithTag",state.AI.position,16,"citizen")[1]
				local civ = query( "gameSpatialDictionary","gridGetCivilization",state.AI.position)[1]
				
				if not citizenNearbyResults and civ > 0 then
					send("rendCommandManager",
						"odinRendererCreateParticleSystemMessage",
						"QuagSmokePuff",
						state.AI.position.x,
						state.AI.position.y)
					
					send(SELF,"despawn")
					return
				end
			end
          end
    
		if state.AI.bools["dead"] then
			if not SELF.tags["buried"] then
				send(SELF, "corpseUpdate")
			end
			return
		end
		
		if SELF.tags.friendly_agent then
			local isDay = query("gameSession","getSessionBool","isDay")[1]
			if isDay then
				send("gameSpatialDictionary","gridExploreFogOfWar",
					state.AI.position.x, state.AI.position.y,
					state.AI.ints["sightRadius"])
			else
				-- isNight
				send("gameSpatialDictionary", "gridExploreFogOfWar",
					state.AI.position.x, state.AI.position.y,
					math.ceil(state.AI.ints["sightRadius"] * 0.75) )
			end
		end

		if state.AI.curJobInstance == nil then
               state.AI.canTestForInterrupts = true -- reset testing for interrupts 
			local results = query("gameBlackboard",
							  "gameAgentNeedsJobMessage",
							  state.AI,
							  SELF )

			if results.name == "gameAgentAssignedJobMessage" then 
				state.AI.curJobInstance = results[1]
				state.AI.curJobInstance.assignedCitizen = SELF
				
				send("rendOdinCharacterClassHandler",
					"odinRendererSetCharacterAttributeMessage",
					state.renderHandle,
					"currentJob",
					state.AI.curJobInstance.displayName)	

			end
          else
               -- interrupt only at 1s intervals because enemies getting stuff right isn't as important as humans
               local oneSecond = false
               if state.AI.ints.updateTimer % 10 == 0 then
                    oneSecond = true
               end
               
               if oneSecond then
                    local results = query( "gameBlackboard",
								  "gameAgentTestForInterruptsMessage",
								  state.AI,
								  SELF )
				
                    if results.name == "gameAgentAssignedJobMessage" then
                         results[1].assignedCitizen = SELF
                         if state.AI.curJobInstance then
                              if state.AI.FSMindex > 0 then
                                   -- run the abort state
                                   local tag = state.AI.curJobInstance:findTag( state.AI.curJobInstance:FSMAtIndex( state.AI.FSMindex  ) )
                                   local name = state.AI.curJobInstance:findName ( state.AI.curJobInstance:FSMAtIndex( state.AI.FSMindex ) )
							
							--printl("DAVID",state.AI.name .. " aborting " .. name )
								  
                                   -- load up our fsm 
                                   local FSMRef = state.AI.curJobInstance:FSMAtIndex(state.AI.FSMindex)
                                   
                                   local targetFSM
                                   if FSMRef:isFSMDisabled() then
                                        targetFSM = ErrorFSMs[ FSMRef:getErrorFSM() ]
                                   else 
                                        targetFSM = FSMs[ FSMRef:getFSM() ]
                                   end
     
                                   local ok
                                   local nextState
                                   ok, errorState = pcall( function() targetFSM[ "abort" ](state, tag, name) end )
          
                                   if not ok then 
                                        print("ERROR: " .. errorState )
                                        FSM.stateError( state )
                                   end
                              end
     
                              state.AI.curJobInstance:abort( "Interrupt hit." )
                              state.AI.curJobInstance = nil					
                         end
     
                         state.AI.abortJob = true
                         if reason ~= nil and reason ~= "" then
                              state.AI.abortJobMessage = reason 
                         end 
                         if state.AI.abortJobMessage == nil then
                              state.AI.abortJobMessage = ""
                         end
     
                         -- Reset the counter for next time
                         state.AI.FSMindex = 0
                         state.AI.curJobInstance = results[1]
                         send("rendOdinCharacterClassHandler",
						"odinRendererSetCharacterAttributeMessage",
                              state.renderHandle,
                              "currentJob",
                              state.AI.curJobInstance.displayName)
					
                    else
                         -- slightly awkward, but we need to do this on the one second tick too.
                         local keepStepping = true
                         while keepStepping do
                              keepStepping = FSM.step(state) 
                         end	
                    end
			else
				local keepStepping = true
				while keepStepping do
					keepStepping = FSM.step(state) 
				end		
			end
          end
	>>
  
	respond getHomeLocation()
	<<
		if SELF.tags["idle"] then
			state.wanderDestination.x = state.spawnLocation.x + rand(-10, 10)
			state.wanderDestination.y = state.spawnLocation.y + rand(-10, 10)
			return "getHomeLocationResponse", state.wanderDestination
		else
			return "getHomeLocationResponse", state.spawnLocation
		end
	>>

	receive despawn() override
	<<
		if SELF.tags["buriedandhidden"] then
			FSM.abort( state, "Despawning.")
			if state.AI.possessedObjects then
				local holdingChar = false
				for key,value in pairs(state.AI.possessedObjects) do
					if key == "curPickedUpCharacter" then
						send("rendOdinCharacterClassHandler", "odinRendererCharacterDetachCharacter", state.renderHandle, value.id, "Bones_Group");
						send(value, "DropItemMessage", state.AI.position.x, state.AI.position.y)
						send(value, "GameObjectPlace", state.AI.position.x, state.AI.position.y)
						send("rendOdinCharacterClassHandler",
							"odinRendererSetCharacterAnimationMessage",
							value.id,
							"corpse_dropped", false)
						
						holdingChar = true
					elseif key then
						send(value, "DestroySelf", state.AI.curJobInstance )
					end
				end
			end
			
			if state.AI.possessedObjects["curCarriedTool"] then
				send(state.AI.possessedObjects["curCarriedTool"], "DestroySelf", state.AI.curJobInstance )
			end
			
			if state.AI.possessedObjects["curPickedUpItem"] then
				send(state.AI.possessedObjects["curPickedUpItem"], "DestroySelf", state.AI.curJobInstance )
			end
			send(SELF,"AICancelJob", "despawning")
			send(SELF,"ForceDropEverything")
			send("rendUIManager", "uiRemoveColonist", SELF.id)
			send("gameSpatialDictionary", "gridRemoveObject", SELF)
			send("rendOdinCharacterClassHandler", "odinRendererDeleteCharacterMessage", state.renderHandle)
			send("gameBlackboard", "gameObjectRemoveTargetingJobs", SELF, nil)
			destroyfromjob(SELF, ji)
		else
			-- disappear in a poof of smoke.
			send("rendCommandManager",
				"odinRendererCreateParticleSystemMessage",
				"CeramicsStamperPoof",
				state.AI.position.x,
				state.AI.position.y)
		
			-- remove me from world, take carried objects with me.
			-- if carrying a body (for some reason), leave it behind.
			
			if state.AI.possessedObjects then
				local holdingChar = false
				for key,value in pairs(state.AI.possessedObjects) do
					if key == "curPickedUpCharacter" then
						send("rendOdinCharacterClassHandler", "odinRendererCharacterDetachCharacter", state.renderHandle, value.id, "Bones_Group");
						send(value, "DropItemMessage", state.AI.position.x, state.AI.position.y)
						send(value, "GameObjectPlace", state.AI.position.x, state.AI.position.y)
						send("rendOdinCharacterClassHandler",
							"odinRendererSetCharacterAnimationMessage",
							value.id,
							"corpse_dropped", false)
						
						holdingChar = true
					elseif key then
						send(value, "DestroySelf", state.AI.curJobInstance )
					end
				end
			end
			
			if state.AI.possessedObjects["curCarriedTool"] then
				send(state.AI.possessedObjects["curCarriedTool"], "DestroySelf", state.AI.curJobInstance )
			end
			
			if state.AI.possessedObjects["curPickedUpItem"] then
				send(state.AI.possessedObjects["curPickedUpItem"], "DestroySelf", state.AI.curJobInstance )
			end
			
			send("rendOdinCharacterClassHandler", "odinRendererDeleteCharacterMessage", state.renderHandle)
			send("gameSpatialDictionary", "gridRemoveObject", SELF)
			destroyfromjob(SELF,ji)
		end
	>>

	receive deathBy( gameObjectHandle damagingObject, string damageType )
	<<
		-- TODO flesh out handling of damagingObject and damageType into interesting descriptions.

		SELF.tags.meat_source = true
		
		-- explode into meats if blown up
		if damageType == "explosion" or damageType == "shrapnel" then
			meat_splosion()
		end

		local animName = bipedDeathAnimSmart(damageType) -- func in ai_agent.go
		
		if animName then
			send("rendOdinCharacterClassHandler",
				"odinRendererSetCharacterAnimationMessage",
				state.renderHandle,
				animName,
				false)
		end
    
		if state.AI.curJobInstance then
			FSM.abort( state, "Died.")
		end
		
		if SELF.tags.occult_inspector then
			send("rendOdinCharacterClassHandler",
				"odinRendererSetCharacterCustomTooltipMessage",
				SELF.id,
				"ui\\tooltips\\occultInspectorDeadTooltip.xml")
		else
			send("rendOdinCharacterClassHandler",
				"odinRendererSetCharacterCustomTooltipMessage",
				SELF.id,
				"ui\\tooltips\\clockworkianFriendlyDeadTooltip.xml")
		end
		
		send(SELF,"resetInteractions")
	>>

	receive HarvestMessage( gameObjectHandle harvester, gameSimJobInstanceHandle ji )
	<<
		local numSteaks = 2 -- so gross.

		for s=1, numSteaks do
			
			local r = rand(1,2)
			if r == 1 then
				if SELF.tags["armoured"] then
					results = query("scriptManager",
								 "scriptCreateGameObjectRequest",
								 "item",
								 {legacyString = "iron_plates"} )
				else
					results = query("scriptManager",
								 "scriptCreateGameObjectRequest",
								 "item",
								 {legacyString = "long_pork"} )
				end
			else
				results = query("scriptManager",
							 "scriptCreateGameObjectRequest",
							 "item",
							 {legacyString = "long_pork"} )
			end
			
			handle = results[1];
			 if( handle == nil ) then 
				printl("Creation failed")
				return "abort"
			 else 
				local range = 1
				positionResult = query("gameSpatialDictionary", "nearbyEmptyGridSquare", state.AI.position, range)
				while not positionResult[1] do
					range = range + 1
					positionResult = query("gameSpatialDictionary", "nearbyEmptyGridSquare", state.AI.position, range)
				end
				if positionResult[1].onGrid then
					send( handle, "GameObjectPlace", positionResult[1].x, positionResult[1].y  )
				else
					send( handle, "GameObjectPlace", state.AI.position.x, state.AI.position.y  )
				end
			end
		end
		
		state.AI.bools["rotted"] = true
		SELF.tags["meat_source"] = nil
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterGeometry", 
			state.renderHandle,
			"models\\character\\body\\bipedSkeleton.upm", 
			"models\\character\\heads\\headSkull.upm",
			"none",
			"biped",
			"idle_dead")
	>>
  
     respond getIdleAnimQueryRequest()
	<<
		-- return a random animation from predefined in characters.xml
		-- this is kinda hardcoding stuff, sadly.
		local animName = "idle" -- idle is class-based so we play it more often
		if state.AI.ints["health"] < state.AI.ints["healthMax"] - 2 then
			animName = "idle_injured"
		else
			local idleAnims = {"idle", "idle_alt1", "idle_alt2", "idle_alt3",}
			animName = idleAnims[rand(1,#idleAnims)]
		end

		return "idleAnimQueryResponse", animName
	>>
     
	receive spawnGibs()
	<<  
		if state.entityData["gibs"] then
			for s=1, rand( state.entityData["gibs"].min , state.entityData["gibs"].max ) do
                    local gibName = state.entityData["gibs"].name
				results = query( "scriptManager", "scriptCreateGameObjectRequest", "clearable", { legacyString = gibName } )
				handle = results[1]
				
                    if not handle then 
                         printl("Creation failed")
                         return "abort"
				else 
                         local range = 1
                              positionResult = query("gameSpatialDictionary", "nearbyEmptyGridSquare", state.AI.position, range)
                              while not positionResult[1] do
                                   range = range + 1
                                   positionResult = query("gameSpatialDictionary", "nearbyEmptyGridSquare", state.AI.position, range)
                              end
                         
                         if positionResult[1].onGrid then
                              send( handle, "GameObjectPlace", positionResult[1].x, positionResult[1].y  )
                         else
                              send( handle, "GameObjectPlace", state.AI.position.x, state.AI.position.y  )
                         end
			    end
			end
		else
               -- no gibs :(
		end
	>>
     
	receive corpseUpdate()
	<<
		if not SELF.tags["corpse_interact"] then
			send(SELF, "resetInteractions")
			SELF.tags["corpse_interact"] = true
		end

		if state.AI.bools["rotted"] then

		else
			--broadcast that there's a rotting corpse over here.
			state.AI.ints["corpse_timer"] = state.AI.ints["corpse_timer"] - 1
			
			if state.AI.ints["corpse_timer"] % 100 == 0 then
				
				if state.AI.ints["corpse_timer"] < state.AI.ints["corpse_vermin_spawn_time_start"] and
					state.numVerminSpawned < 8 then
					
					if rand(1,8) == 8 then
						local handle = query( "scriptManager",
								"scriptCreateGameObjectRequest",
								"vermin",
								{legacyString = "Tiny Beetle" } )[1]
						
						local positionResult = query("gameSpatialDictionary", "nearbyEmptyGridSquare", state.AI.position, 3)[1]
						send(handle,
							"GameObjectPlace",
							positionResult.x,
							positionResult.y  )
						
						state.numVerminSpawned = state.numVerminSpawned +1
					end
				end
			end

			if state.AI.ints["corpse_timer"] <= 0 then
				-- here's your skeleton model swap
				state.AI.bools["rotted"] = true
				state.AI.bools["onFire"] = false -- because we're done with that.
				SELF.tags["burning"] = false
				SELF.tags["meat_source"] = false
				send( "rendOdinCharacterClassHandler",
						"odinRendererSetCharacterGeometry", 
						state.renderHandle,
						"models\\character\\body\\bipedSkeleton.upm", 
						"models\\character\\heads\\headSkull.upm",
						"none",
						"biped",
						"idle_dead")
                    
				send("rendCommandManager",
						"odinRendererCreateParticleSystemMessage",
						"MiasmaBurst",
						state.AI.position.x,
						state.AI.position.y)
                    
			end
		end
	>>

     receive hearExclamation( string name, gameObjectHandle exclaimer, gameObjectHandle firstIgnored )
	<<
		if SELF.tags.dead then return end
		
		if name == "detectCombat" then
			-- decrease morale because combat is scary
			state.AI.ints["morale"] = state.AI.ints["morale"] - 2
			if state.AI.ints["morale"] <= 50 and not SELF.tags["fleeing"] then
				--FSM.abort(state,"Morale broken.")
				SELF.tags["fleeing"] = true
				SELF.tags["coward"] = true
				if state.group then
					send(state.group, "memberMoraleBroken")
				end
			end
		end
	>>
	
     receive damageMessage( gameObjectHandle damagingObject, string damageType, int damageAmount, string onhit_effect )
	<<
		if SELF.tags.dead then return end
		
          -- decrease morale because getting hurt is scary
          state.AI.ints["morale"] = state.AI.ints["morale"] - damageAmount * 2
          if state.AI.ints["morale"] <= 50 and not SELF.tags["fleeing"] then
               FSM.abort(state,"Morale broken." )
               SELF.tags["fleeing"] = true
			SELF.tags["coward"] = true
          end
	>>
	
	receive emoteAffliction()
	<<
	     -- decrease morale after affliction because getting hurt is scary
          state.AI.ints["morale"] = state.AI.ints["morale"] - 3
          if state.AI.ints["morale"] <= 50 and not SELF.tags["fleeing"] then
               FSM.abort(state,"Morale broken." )
               SELF.tags["fleeing"] = true
			SELF.tags["coward"] = true
			if state.group then
				send(state.group, "memberMoraleBroken")
			end
          end
	>>
     
     receive resetEmoteTimer()
	<<
		state.AI.ints["emoteTimer"] = 0
	>>
	
	receive convertToCitizen()
	<<
		-- TODO copy from bandit if we do this.
	>>

	receive incMorale(int morale)
	<<
		state.AI.ints.morale = state.AI.ints.morale + morale
	>>
	
	receive resetInteractions()
	<<
		printl("ai_agent", state.AI.name .. " received resetInteractions")
		
		send("rendInteractiveObjectClassHandler",
				"odinRendererClearInteractions",
				state.renderHandle)
		
		if SELF.tags.dead and
			not state.assignment then
			
			send("rendInteractiveObjectClassHandler",
                    "odinRendererAddInteractions",
				state.renderHandle,
                         "Give " .. state.AI.name .. " a Proper Burial",
                         "Bury Corpse (player order)",
                         "Bury Corpses", --"Bury Corpses",
                         "Bury Corpse (player order)", --"Bury Corpse (player order)",
						"graveyard",
						"",
						"Dirt",
						false,true)
			
			send("rendInteractiveObjectClassHandler",
                    "odinRendererAddInteractions",
				state.renderHandle,
                         "Dump the Corpse of " .. state.AI.name,
                         "Dump Corpse (player order)",
                         "Dump Corpses", --"Dump Corpses",
                         "Dump Corpse (player order)", --"Dump Corpse (player order)",
						"graveyard",
						"",
						"Dirt",
						false,true)
			
		end
	>>
	
	receive InteractiveMessage( string messagereceived )
	<<
		send(SELF,"HandleInteractiveMessage",messagereceived,nil)
	>>
     
     receive InteractiveMessageWithAssignment( string messagereceived, gameSimAssignmentHandle assignment )
     <<
		send(SELF,"HandleInteractiveMessage",messagereceived,assignment)
     >>
	
	receive HandleInteractiveMessage(string messagereceived, gameSimAssignmentHandle assignment)
	<<
		printl("ai_agent", state.AI.name .. " receive HandleInteractiveMessage: " .. messagereceived)
		local setCancelInteraction = false
		
		if messagereceived == "Bury Corpse (player order)" and
			not state.assignment then
			
			send("gameBlackboard",
				"gameObjectRemoveTargetingJobs",
				SELF,
				nil)
			
			send("rendInteractiveObjectClassHandler",
				"odinRendererClearInteractions",
				state.renderHandle)
			
			send("rendCommandManager",
				"odinRendererCreateParticleSystemMessage",
				"Small Beacon",
				state.AI.position.x,
				state.AI.position.y)
			
			if not assignment then
				assignment = query("gameBlackboard",
								"gameObjectNewAssignmentMessage",
								SELF,
								"Burial",
								"",
								"")[1]
			end
			
			send("rendOdinCharacterClassHandler",
				"odinRendererCharacterExpression",
				state.renderHandle,
				"dialogue",
				"jobshovel",
				true)
			
			send( "gameBlackboard",
				"gameObjectNewJobToAssignment",
				assignment,
				SELF,
				"Bury Corpse (player order)",
				"body",
				true )
			
			setCancelInteraction = true
			state.assignment = assignment
			
		elseif messagereceived == "Dump Corpse (player order)" and
			not state.assignment then
			
			send("gameBlackboard",
				"gameObjectRemoveTargetingJobs",
				SELF,
				nil)
			
			send("rendInteractiveObjectClassHandler",
				"odinRendererClearInteractions",
				state.renderHandle)

			send("rendCommandManager",
				"odinRendererCreateParticleSystemMessage",
				"Small Beacon",
				state.AI.position.x,
				state.AI.position.y)
			
			if not assignment then
				assignment = query("gameBlackboard",
								"gameObjectNewAssignmentMessage",
								SELF,
								"Dump Corpse",
								"",
								"")[1]
			end
			
			send("gameBlackboard",
				"gameObjectNewJobToAssignment",
				assignment,
				SELF,
				"Dump Corpse (player order)",
				"body",
				true )
			
			send("rendOdinCharacterClassHandler",
				"odinRendererCharacterExpression",
				state.renderHandle,
				"dialogue",
				"jobhand",
				true)
			
			setCancelInteraction = true
			state.assignment = assignment
			
		elseif messagereceived == "Cancel corpse orders" then
			send("gameBlackboard",
				"gameObjectRemoveTargetingJobs",
				SELF,
				nil)
			
			send(SELF,"resetInteractions")
			state.assignment = nil
		end
		
		if setCancelInteraction then
			send("rendInteractiveObjectClassHandler",
                    "odinRendererAddInteractions",
				state.renderHandle,
                         "Cancel orders for corpse of " .. state.AI.name,
                         "Cancel corpse orders",
                         "Cancel corpse orders", --"Cancel corpse orders",
                         "Cancel corpse orders", --"Cancel corpse orders",
						"graveyard",
						"",
						"Dirt",
						false)
		end
	>>
	
	receive AssignmentCancelledMessage( gameSimAssignmentHandle assignment )
	<<
		printl("ai_agent", state.AI.name .. " received AssignmentCancelledMessage")
		send("rendInteractiveObjectClassHandler",
			"odinRendererRemoveInteraction",
			state.renderHandle,
			"Cancel Assignment")
          
		state.assignment = nil
		send(SELF,"resetInteractions")
	>>

	receive JobCancelledMessage(gameSimJobInstanceHandle job)
	<<
		printl("ai_agent", state.AI.name .. " received JobCancelledMessage")
		state.assignment = nil
		send(SELF,"resetInteractions")
	>>
>>