gameobject "foreigner" inherit "ai_agent"
<<
	local
	<<
          function foreigner_doOneSecondUpdate()
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
			printl("foreigner", "foreigner type not found")
			return
		end
          state.entityData = entityData
	    
          state.AI.strs["mood"] = "foreign"
		
		local nationInfo = EntityDB[ entityData.nation .. "Info"]
		state.nation = entityData.nation
          
          if rand(0,100) < 50 then 
			state.AI.strs["gender"] = "male"
               state.AI.strs["firstName"] =  nationInfo.maleNames[ rand(1, #nationInfo.maleNames ) ]
		else
			state.AI.strs["gender"] = "female"
               state.AI.strs["firstName"] =  nationInfo.femaleNames[ rand(1, #nationInfo.femaleNames ) ]
		end
		
		if nationInfo.name == "NovorusInfo" then
			-- gendered family names for Novorusians. (Thx Kristina.)
			if state.AI.strs.gender == "female" then
				state.AI.strs["lastName"] = nationInfo.femaleLastNames[ rand(1, #nationInfo.femaleLastNames ) ]
			else
				state.AI.strs["lastName"] = nationInfo.maleLastNames[ rand(1, #nationInfo.maleLastNames ) ]
			end
		else
			state.AI.strs["lastName"] = nationInfo.lastNames[ rand(1, #nationInfo.lastNames ) ]
		end
		
		state.AI.name = state.AI.strs.firstName .. " " .. state.AI.strs.lastName
		
		local isArmoured = false
		local variants = {[1]="A", [2]="B", [3]="C", [4]="D"}
		state.AI.strs["variant"] = variants[rand(1,4)]
		
		state.AI.strs["citizenClass"] = entityName
		local models = getModelsForClass( state.AI.strs["citizenClass"],
									state.AI.strs["gender"],
									state.AI.strs["variant"] )

		state.animSet = models["animationSet"]

          send( "rendOdinCharacterClassHandler",
                    "odinRendererCreateCitizen", 
                    SELF, 
                    models["torsoModel"], 
                    models["headModel"],
				"", --models["hairModel"], -- "models/hats/phrygiancap.upm", --models["hairModel"], 
                    entityData.uniformHatModel, --models["hatModel"], 
                    models["animationSet"], 0, 0 )
		
		state.headModel = models["headModel"]
		state.hairModel = models["hairModel"]
       

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
          state.AI.ints["healthMax"] = entityData.health
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

          send("rendOdinCharacterClassHandler", "odinRendererSetCharacterWalkTicks",  SELF.id, state.AI.walkTicks )

		send("gameSpatialDictionary",
			"registerSpatialMapString",
			SELF,
			entityData.occupancyMap,
			entityData.occupancyMapRotate45,
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
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterCustomTooltipMessage",
			SELF.id,
			"ui\\tooltips\\banditTooltip.xml")
          
		-- set up foreign faction flags as appropriate.
		if state.nation == "Stahlmark" then
			send("gameSpatialDictionary", "gameObjectAddBit", SELF, 6)
		elseif state.nation == "Republique" then
			send("gameSpatialDictionary", "gameObjectAddBit", SELF, 5)
		elseif state.nation == "Novorus" then
			send("gameSpatialDictionary", "gameObjectAddBit", SELF, 7)
		end
		
		local myNation = query("gameSession","getSessiongOH",state.nation)[1]
		 
		if query(myNation, "isHostile")[1] and SELF.tags.soldier then
			send(SELF,"makeHostile")
		elseif query(myNation, "isFriendly")[1] then
			send(SELF,"makeFriendly")
		else
			send(SELF,"makeNeutral")
		end
		
		state.exit_timer = 60 * 10 * 4 -- half a day
		
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
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterCustomTooltipMessage",
			SELF.id,
			"ui\\tooltips\\foreignerHostileTooltip.xml")
	>>
	
	receive makeFriendly()
	<<
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle,
			"name",
			state.AI.name)
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterCustomTooltipMessage",
			SELF.id,
			"ui\\tooltips\\foreignerFriendlyTooltip.xml")
	>>
	
	receive makeNeutral()
	<<
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle,
			"name",
			state.AI.name)
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterCustomTooltipMessage",
			SELF.id,
			"ui\\tooltips\\foreignerTooltip.xml")
	>>
	
	receive SleepMessage()
	<<
		state.asleep = true
	>>

     receive Update()
     <<
  		if state.AI.thinkLocked then
               return
          end
		
          if state.AI.ints.updateTimer % 10 == 0 then
               foreigner_doOneSecondUpdate()
          end
    
		if state.AI.ints.updateTimer % 29 == 0 then
			SELF.tags.helping_friend = nil
		end
    
		if state.AI.bools["dead"] then
			if not SELF.tags["buried"] then
				send(SELF, "corpseUpdate")
			end
			return
		end
		
		if SELF.tags.exit_map then
			
			if state.exit_timer == 0 then
				-- force despawn
				send(SELF,"despawn")
				return
			else
				state.exit_timer = state.exit_timer - 1
			end
		end
		
		if SELF.tags.friendly_agent then
			local isDay = query("gameSession","getSessionBool","isDay")[1]
			if isDay then
				send("gameSpatialDictionary","gridExploreFogOfWar",
					state.AI.position.x,
					state.AI.position.y,
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
	>>

	receive deathBy( gameObjectHandle damagingObject, string damageType )
	<<
		-- TODO flesh out handling of damagingObject and damageType into interesting descriptions.

		SELF.tags.meat_source = true

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
		
		local civ = query("gameSpatialDictionary",
						"gridGetCivilization",
						state.AI.position)[1]
						
		if civ == 0 then
			send(SELF,"HandleInteractiveMessage","Bury Corpse (player order)",nil)
		else
			send(SELF,"resetInteractions")
		end
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
		
		send(SELF,"resetInteractions")
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
		-- TODO set something up for corpses in general.
		--[[elseif name == "detectBanditCorpse" then

			state.AI.ints["morale"] = state.AI.ints["morale"] - 3
			if state.AI.ints["morale"] <= 50 and not SELF.tags["fleeing"] then
				--FSM.abort(state,"Morale broken.")
				SELF.tags["fleeing"] = true
				SELF.tags["coward"] = true
				if state.group then
					send(state.group, "memberMoraleBroken")
				end
			end	--]]
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
			if state.group then
				send(state.group, "memberMoraleBroken")
			end
          end
		
		if damagingObject then 
			local attackerTags = query(damagingObject, "getTags")[1]
			
			if attackerTags.animal or
				attackerTags.foreigner or
				attackerTags.obeliskian or
				attackerTags.bandit then
				
				if SELF.tags.friendly_agent then
					send(damagingObject,"addTag","hostile_agent")
				end
				
				if state.group and damagingObject then
					send(state.group,"helpMemberInCombat", SELF, damagingObject)
				end
			end
		end
		
		if SELF.tags.middle_class and SELF.tags.trader and (state.AI.ints.health < 5 or state.AI.ints.morale <50) then
			-- you're the group leader.
			-- if you take a bunch of damage, then screw this.
			if not SELF.tags.exit_map_flee and not SELF.tags.notified_of_trader_flight then
				if state.group then
					send(state.group,"forceRetreat")
					
					SELF.tags.notified_of_trader_flight = true
					--local gname = query(state.group,"getName")[1]
					--
					send("rendCommandManager",
						"odinRendererStubMessage",
						"ui\\thoughtIcons.xml", -- iconskin
						"retreat", -- icon
						"Traders flee due to attacks!", -- header text
						"A group of traders was forced to flee from your colony due to being attacked. They don't want to trade in such an unsafe place.", -- text description
						"Right-click to dismiss.", -- action string
						"tradersFlee", -- alert type (for stacking)
						"", -- imagename for bg
						"low", -- importance: low / high / critical -- high
						nil, -- object ID
						30 * 1000, -- duration in ms
						0, -- "snooze" time if triggered multiple times in rapid succession
						SELF)
				end
			end
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
	
	receive helpFriendInCombat(gameObjectHandle victim, gameObjectHandle attacker)
	<<
		if SELF.tags.helping_friend then return end
		
		local importantJobs = {
			"Flee (agent, from horror)",
			"Shoot at Target (neutral_agent)",
			"Melee Attack With Weapon (neutral_agent)",
			"Shoot at Target (hostile_agent)",
			"Melee Attack With Weapon (hostile_agent)",
			"Shoot at Target (friendly_agent)",
			"Melee Attack With Weapon (friendly_agent)",
		}
		
		if state.AI.curJobInstance then 
			for k,v in pairs(importantJobs) do
				if v == state.AI.curJobInstance.name then
					-- don't interrupt.
					return
				end
			end
		end
		
		-- 1. am I close enough to help?
		local otherPos = query(victim, "gridGetPosition")[1]
		
		local xdist = math.abs( state.AI.position.x - otherPos.x )
		local ydist = math.abs( state.AI.position.y - otherPos.y )
		local dist = math.floor( math.sqrt( xdist*xdist + ydist*ydist ) )

		if dist <= 12 then
			-- attack
			FSM.abort(state, "need to assist friend")
			SELF.tags.helping_friend = true
		else
			-- too far.
			return
		end
		
		-- 2. am I really, really busy?
		-- uh, TODO.
		
		-- If none of the above, let's attack.
		
		-- ranged or melee?
		if SELF.tags.has_ranged_attack then
			send( "gameBlackboard",
				"gameCitizenJobToMailboxMessage",
				SELF,
				attacker,
				"Shoot at Target (neutral_agent)",
				"target")
			
		elseif SELF.tags.has_melee_attack and
			not (SELF.tags.noncombatant or SELF.tags.trader or SELF.tags.civilian) then
			
			send( "gameBlackboard",
				"gameCitizenJobToMailboxMessage",
				SELF,
				attacker,
				"Melee Attack With Weapon (neutral_agent)",
				"target")
		end
		
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
                         "", --"Bury Corpses",
                         "", --"Bury Corpse (player order)",
						"graveyard",
						"",
						"Dirt",
						false,true)
			
			send("rendInteractiveObjectClassHandler",
                    "odinRendererAddInteractions",
				state.renderHandle,
                         "Dump the Corpse of " .. state.AI.name,
                         "Dump Corpse (player order)",
                         "Dump Corpses",
                         "Dump Corpse (player order)",
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
                         "Cancel corpse orders",
                         "Cancel corpse orders",
						"graveyard",
						"",
						"Dirt",
						false, true)
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