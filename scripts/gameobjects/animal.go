gameobject "animal" inherit "ai_agent"
<<
	local
	<<
          -- This update occurs once per game second.
          function animal_doSecondUpdate()
               --state.AI.ints["hunger"] = state.AI.ints["hunger"] + 1
               state.AI.ints["emoteTimer"] = state.AI.ints["emoteTimer"] + 1
               
               if state.AI.ints["morale"] < 100 then
                    if state.AI.ints["morale"] < 0 then
                         state.AI.ints["morale"] = 0
                    end
                    state.AI.ints["morale"] = state.AI.ints["morale"] + 1
               end
               -- check injuries
               animal_resetHealthTags()
               
               -- do health regen
               if state.AI.ints["healthTimer"] ~= -1 then
                    if state.AI.ints["health"] < state.AI.ints["healthMax"] then
                         state.AI.ints["healthTimer"] = state.AI.ints["healthTimer"] - 1
                         if state.AI.ints["healthTimer"] == 0 then
                              state.AI.ints["health"] = state.AI.ints["health"] + 1
                              state.AI.ints["healthTimer"] = state.AI.ints["healthTimerMax"]
                         end
                    end
               end
			
			-- do hunger check
			if not state.AI.ints.secondsTimer then
				state.AI.ints.secondsTimer = 1
			else
				state.AI.ints.secondsTimer = state.AI.ints.secondsTimer + 1
			
				if state.AI.ints.secondsTimer > EntityDB.WorldStats.dayNightCycleSeconds then
					state.AI.ints.secondsTimer = 0
					if state.AI.ints.hunger < 4 then
						state.AI.ints.hunger = state.AI.ints.hunger + 1
					end
				end
			end
			
			if query("gameSpatialDictionary",
					"gridHasSpatialTag",
					state.AI.position,
					"occupiedByStructure")[1] then
				
				if not state.buildingTimer then
					state.buildingTimer = 0
				end
				
				state.buildingTimer = state.buildingTimer + 1
				
				-- find a loc that isn't inside!
				-- drop outside foundation so we don't get floaters.

				if state.buildingTimer > 10 then 
					local newLoc = gameGridPosition:new()
					local x = state.AI.position.x
					local y = state.AI.position.y
					local isInvalidDrop = true
					local i = 1
					while isInvalidDrop do
						
						newLoc.x = x + rand(i * -1,i)
						newLoc.y = y + rand(i * -1,i)
						
						isInvalidDrop = query("gameSpatialDictionary","gridHasSpatialTag",newLoc,"occupiedByStructure" )[1]
						i = i + 1
					end
					
					send(SELF, "GameObjectPlace", newLoc.x, newLoc.y)
				end
			else
				state.buildingTimer = 0
			end
          end
          
          -- This update occurs once per 3 game seconds.
          function animal_doThreeSecondUpdate()
			
          end


          function animal_resetHealthTags()
               if state.AI.ints["health"] > 0 then
                    -- if injured at all, "injured"
                    -- if injured > 50%, "badly_injured"
                    if state.AI.ints["health"] >= state.AI.ints["healthMax"] then
                         SELF.tags["injured"] = nil
                         SELF.tags["badly_injured"] = nil
                    elseif state.AI.ints["health"] < state.AI.ints["healthMax"] then
                         SELF.tags["injured"] = true
                    end
                    
                    if state.AI.ints["health"] <= div(state.AI.ints["healthMax"],2) then
                         SELF.tags["badly_injured"] = true
                    else
                         SELF.tags["badly_injured"] = nil
                    end
               else
                    -- deaded.
                    SELF.tags["injured"] = nil
                    SELF.tags["badly_injured"] = nil
               end
          end
     
          function setposition( x, y )
               local newPos = gameGridPosition:new()
               newPos:set( x, y )
               state.AI.position = newPos
          end
	>>

	state
	<<
		gameAIAttributes AI
		string animalclass
		string animSet
          string entityName
          table entityData
		int renderHandle
	>>

	receive Create( stringstringMapHandle init )
	<<
          state.entityName = init["legacyString"]
          state.entityData = EntityDB[ state.entityName ]
		local entityData = state.entityData
		if entityData == nil then
			printl("animal", "Entitydata nil for " .. state.entityName )
		end
		
		local worldstats = EntityDB["WorldStats"]
		
		state.AI.strs["firstName"] = state.entityName 
		state.AI.strs["lastName"] = tostring (SELF.id)
		
		state.AI.name = state.AI.strs.firstName .. " " .. state.AI.strs.lastName
          
		state.AI.strs["citizenClass"] = init["legacyString"]
          state.AI.strs["name"] = state.AI.strs["citizenClass"]
		state.AI.ints["num_afflictions"] = 0
          state.AI.ints["max_afflictions"] = 1
		state.AI.ints["health"] = 3
		state.AI.ints["healthMax"] = 3
		state.AI.ints["healthTimer"] = -1 -- never regenerate! Requires repairs. 
		state.AI.ints["health"] = state.AI.ints["healthMax"]
          state.AI.ints["morale"] = 100
          state.AI.ints["hunger"] = rand(0, 2 )
		state.AI.ints["tiredness"] = rand(0, 2 )

		setposition( 0, 0 )
          
		local modelM = entityData.modelM
		local modelF = entityData.modelF
		
          if rand(1,2) == 1 then
               send("rendOdinCharacterClassHandler",
				"odinRendererCreateCharacter", 
                    SELF,
				modelM,
				entityData.animationSetM,
				0,
				0 )
			
               state.AI.strs["gender"] = "male"
               state.animSet = entityData.animationSetM
          else
               send("rendOdinCharacterClassHandler",
				"odinRendererCreateCharacter", 
                    SELF,
				modelF,
				entityData.animationSetF,
				0,
				0 )
			
               state.AI.strs["gender"] = "female"
               if entityData.animationSetF then
                    state.animSet = entityData.animationSetF
               else
                    state.animSet = entityData.animationSetM
               end
          end

		SELF.tags = {}
		if entityData.job_classes then
			for k,v in pairs(entityData.job_classes) do
				SELF.tags[v] = true
			end
		end
          
          if entityData.tags then
               for k,v in pairs(entityData.tags) do
                    SELF.tags[v] = true
               end
          else
               SELF.tags = { "animal", "huntable" }
          end
		
          state.AI.walkTicks = entityData.walkTicks
		
          send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterWalkTicks", 
               SELF.id,
               state.AI.walkTicks )

		if entityData.health then
               state.AI.ints.healthMax = entityData.health.healthMax
               state.AI.ints.health = entityData.health.healthMax
               state.AI.ints.healthTimer = entityData.health.healthTimerSeconds
               state.AI.ints.healthTimerMax = entityData.health.healthTimerSeconds
          end
		
		-- Register our mass with the gameSpatialDictionary
		send( "gameSpatialDictionary",
                    "registerSpatialMapString",
                    SELF,
                    entityData.occupancyMap,
                    entityData.occupancyMapRotate45,
                    true )
          
		state.AI.locs["herdCentre"] = state.AI.position
		
		if SELF.tags.herbivore then
			send("gameSpatialDictionary", "gameObjectAddBit", SELF, 9)
		elseif SELF.tags.vicious then
			send("gameSpatialDictionary", "gameObjectAddBit", SELF, 8)
		elseif SELF.tags.large_carnivore then
			send("gameSpatialDictionary", "gameObjectAddBit", SELF, 8)
		else
			-- HAX. need case for small non-hostile carnivores. Or do we?
			send("gameSpatialDictionary", "gameObjectAddBit", SELF, 9)
		end

		send("rendOdinCharacterClassHandler",
				"odinRendererSetCharacterAttributeMessage",
				state.renderHandle,
				"animalType",
				state.AI.strs["citizenClass"])
		
--[[
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterCustomTooltipMessage",
			SELF.id,
			"ui\\tooltips\\animalTooltip.xml")
--]]		
		ready()
	>>

	receive GameObjectPlace( int x, int y ) 
	<<
		send(SELF,"resetInteractions")
	>>

	receive Update()
	<<
		if state.asleep then
			--return
		end

		if state.AI.bools["dead"] then
			return
		end
          
          -- get a job, animal!
		if state.AI.curJobInstance == nil then
            state.AI.canTestForInterrupts = true -- reset testing for interrupts 
			local results = query( "gameBlackboard", "gameAgentNeedsJobMessage", state.AI, SELF )

			if results.name == "gameAgentAssignedJobMessage" then 
				state.AI.curJobInstance = results[ 1 ]
				state.AI.curJobInstance.assignedCitizen = SELF
				send("rendOdinCharacterClassHandler",
						"odinRendererSetCharacterAttributeMessage",
						state.renderHandle,
						"currentJob",
						state.AI.curJobInstance.displayName)
			end
          else
               -- NOTE: interrupt only at 1s intervals because animals getting their
			-- jobs correct isn't as important as it is for humans
               local oneSecond = false
               if state.AI.ints.updateTimer % 10 == 0 then
                    oneSecond = true
				
               end
               
               if oneSecond then 
                    local results = query( "gameBlackboard", "gameAgentTestForInterruptsMessage", state.AI, SELF )
                    if results.name == "gameAgentAssignedJobMessage" then
                         results[1].assignedCitizen = SELF
                         if state.AI.curJobInstance then
                              if state.AI.FSMindex > 0 then
                                   -- run the abort state
                                   local tag = state.AI.curJobInstance:findTag( state.AI.curJobInstance:FSMAtIndex( state.AI.FSMindex  ) )
                                   local name = state.AI.curJobInstance:findName ( state.AI.curJobInstance:FSMAtIndex( state.AI.FSMindex ) )
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
                         state.AI.curJobInstance = results[ 1 ]
                         send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage",
                                   state.renderHandle,
                                   "currentJob",
                                   state.AI.curJobInstance.displayName)
                    else
                         -- slightly awkward, but we need to do this on the one second tick too.
                         local keepStepping = true
                         while keepStepping do
                              keepStepping = FSM.step( state ) 
                         end	
                    end
			else
				local keepStepping = true
				while keepStepping do
					keepStepping = FSM.step( state ) 
				end		
			end
		end
          
          if state.AI.ints.updateTimer % 10 == 0 then
               animal_doSecondUpdate()
          end
		
          if state.AI.ints.updateTimer % 30 == 0 then
               animal_doThreeSecondUpdate()
          end
	>>
	
     receive HarvestMessage( gameObjectHandle harvester, gameSimJobInstanceHandle ji )
	<<		
		local entityData = EntityDB[ state.AI.strs["citizenClass"] ] 
		if entityData == nil then
			return
		end
		
		local harvester_tags = query(harvester,"getTags")[1]

          -- BONES
		if entityData.bonesOutput then
			local handle = query("scriptManager",
							"scriptCreateGameObjectRequest",
							"forageSource",
							{legacyString = entityData.bonesOutput } )[1]
			
			if not handle then 
			    printl("animal", "WARNING: bones creation failed")
			    return "abort"
			else 
				local range = 1
				local positionResult = query("gameSpatialDictionary",
									    "nearbyEmptyGridSquare",
									    state.AI.position,
									    range)
				
				while not positionResult[1] do
					range = range + 1
					positionResult = query("gameSpatialDictionary",
									   "nearbyEmptyGridSquare",
									   state.AI.position,
									   range)
					
				end
				
				if positionResult[1].onGrid then
					send( handle, "GameObjectPlace", positionResult[1].x, positionResult[1].y  )
				else
					send( handle, "GameObjectPlace", state.AI.position.x, state.AI.position.y  )
				end
			end
		end
     
		if harvester_tags.citizen then
			
			-- MEATS     
			-- random range output items, if any
			if entityData.butcherOutput then
				if entityData["numCommodityOutput"] then
					numOutput = entityData["numCommodityOutput"]
				else
					numOutput = 2
				end
				
				for s=1, numOutput do

					local handle = query( "scriptManager",
								"scriptCreateGameObjectRequest",
								"item",
								{legacyString = entityData.butcherOutput } )[1]
					
					if not handle then 
					    printl("animal", "WARNING: meat output creation failed")
					    return "abort"
					else
						--local harvesterTags = query(harvester, "getTags")[1]
						local ownerTags = harvester:getOwnerTags()
	
						if harvester_tags.animal then
							send(handle, "ForbidItem")
							local civ = query("gameSpatialDictionary",
											"gridGetCivilization",
											state.AI.position)[1]
						  
							if civ < 10 then
								send(handle,"ClaimItem")
							else
								send(handle,"ForbidItem")
							end
						else
							send(handle, "ClaimItem")
						end
	
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
			end
		end
          
          -- AND MAKE A MESS
          send("rendCommandManager",
			"odinRendererCreateParticleSystemMessage",
			"BloodSplashCentered",
			state.AI.position.x,
			state.AI.position.y)
		
          send(SELF, "spawnGibs")
          
		-- Put some work music in the dynamics.
		incMusic(3,10)

          send("rendOdinCharacterClassHandler", "odinRendererDeleteCharacterMessage", state.renderHandle)
          send("gameSpatialDictionary", "gridRemoveObject", SELF)
          send("gameBlackboard", "gameObjectRemoveTargetingJobs", SELF, ji)

		-- destroy animal?
          SELF.tags["dead_and_harvested"] = true
	>>
	
	receive setHerdCentre(gameGridPosition position)
	<<
		state.AI.locs["herdCentre"] = position
	>>
     
	receive InteractiveMessage( string messagereceived )
	<<
		animal_receiveInteractiveMessage(messagereceived, nil)
	>>
     
     receive InteractiveMessageWithAssignment( string messagereceived, gameSimAssignmentHandle assignment )
     <<
		animal_receiveInteractiveMessage(messagereceived, assignment)
     >>
     
	receive hearExclamation( string name, gameObjectHandle exclaimer, gameObjectHandle firstIgnored )
	<<
		if name == "explosion" then
			send( "gameBlackboard","gameCitizenJobToMailboxMessage", SELF, exclaimer, "Flee From Gunfire", "enemy")
		elseif name == "gunshot" then
			send( "gameBlackboard", "gameCitizenJobToMailboxMessage", SELF, exclaimer, "Flee From Gunfire", "enemy")
		end
	>>

	receive triggerFlee(gameObjectHandle attacker)
	<<
          --[[
          if not SELF.tags["fleeing"] then
               SELF.tags["fleeing"] = true  
               -- decrease morale because combat is scary
               state.AI.ints["morale"] = state.AI.ints["morale"] - 10
               send( "gameBlackboard", "gameCitizenJobToMailboxMessage", SELF, attacker, "Herd Flee", "enemy")
               if VALUE_STORE["showCombatDebugConsole"] then printl("combat", state.AI.name .. " the animal witnessed combat, morale now: " .. state.AI.ints["morale"] ) end
          else
		
          end --]]
	>>
     
	receive damageMessage( gameObjectHandle damagingObject, string damageType, int damageAmount, string onhit_effect )
	<<
          if SELF.tags["rampager"] and not SELF.tags["rampage"] then
             SELF.tags["rampage"] = true  
          end
		
		--[[
		This would be the revised method of handling triggerFlee
		But triggerFlee is commented out, so no point in doing this expensive query.
		
		if SELF.tags["skittish"] then
			
			local results = query("gameSpatialDictionary",
								    "allObjectsInRadiusWithTagRequest",
								    state.AI.position,
								    15,
								    "skittish",
								    true)
			if results then
				if results[1] then
					send(results[1], "triggerFlee", damagingObject)
				end
			end --]]
			
--[[results = query("gameSpatialDictionary", "allObjectsInRadiusRequest", state.AI.position, 20, false)
			if results and results[1] then
				for k,v in pairs(results[1]) do
					if query(v, "getTags")[1] and query(v, "getTags")[1]["skittish"] then
						send(results[1], "triggerFlee", damagingObject)
					end
				end
			end--]]
          --end
     >>

     receive deathBy( gameObjectHandle damagingObject, string damageType )
	<<
		send(SELF,"Vocalize","die")
		
		send("rendOdinCharacterClassHandler",
				"odinRendererSetCharacterAnimationMessage",
				state.renderHandle,
				"death",
				false)

		send("rendOdinCharacterClassHandler",
				"odinRendererQueueCharacterAnimationMessage",
				state.renderHandle,
				"")

		SELF.tags.hostile_entity = nil
		
          if state.entityData.deadTags then
               for k,v in pairs(state.entityData.deadTags) do
                    SELF.tags[v] = true
               end
		end
		SELF.tags.meat_source = true
		 --[[ else
               SELF.tags = {"corpse","meat_source","dead" }
          end--]]

		send("rendInteractiveObjectClassHandler",
			"odinRendererClearInteractions",
			state.renderHandle)
          
          if SELF.tags["meat_source"] then
               send("rendInteractiveObjectClassHandler",
                         "odinRendererAddInteractions",
                         state.renderHandle,
                         "Butcher",
                         "Butcher Corpse (player order)",
                         "", --"Butcher",
                         "", -- "Butcher Corpse (player order)",
                         "corpse",
                         "",
                         "Slice Flesh",
					true,true)
			
			if damagingObject then 
				if damagingObject.tags["citizen"] then
					send( "gameBlackboard",
						"gameCitizenJobToMailboxMessage",
						damagingObject,
						SELF,
						"Butcher Animal Corpse",
						"food")
					
					
					if state.AI.strs["citizenClass"] == "Dodo" then
						send("gameSession","incSessionInt","dodoKillCount", 1)
						local count = query("gameSession","getSessionInt","dodoKillCount")[1]
						if not query("gameSession","getSessionBool","killedManyDodos")[1] and count >= 100 then
							send("gameSession", "setSessionBool", "killedManyDodos", true)
						end
						send("gameSession", "incSteamStat", "stat_dodo_kill_count", 1)
					end
					--[[
					elseif state.AI.strs["citizenClass"] == "Aurochs" then
						send("gameSession","incSessionInt","aurochsKillCount", 1)
						local count = query("gameSession","getSessionInt","aurochsKillCount")[1]
						if not query("gameSession","getSessionBool","killedManyAurochs")[1] and count >= 100 then
							send("gameSession", "setSessionBool", "killedManyAurochs", true)
						end
						send("gameSession", "incSteamStat", "stat_aurochs_kill_count", 1)
						
					end]]
				end
			end
          end

		if state.entityData.spawnAfterDeath then
			local handle = query( "scriptManager",
							"scriptCreateGameObjectRequest",
							state.entityData.spawnAfterDeath.entityType,
							{legacyString = state.entityData.spawnAfterDeath.entityName } )[1]
			
			send( handle,
				"GameObjectPlace",
				state.AI.position.x,
				state.AI.position.y  )
		end
		
		send(SELF,"resetInteractions")
     >>

	receive gameFogOfWarExplored(int x, int y )
	<<
		-- send ping out to all nearby animals
		-- (to wake parts of herd in fog of war)
		local results = query("gameSpatialDictionary",
					 "allObjectsInRadiusWithTagRequest",
					 state.AI.position,
					 6,
					 "animal",
					 true)
		
		-- radius 6 seems to work pretty well.
		if results and results[1] then
			send(results[1], "WakeMessage")
		end
		state.asleep = false
		wake()
	>>
	
	receive ConsumeFood( gameObjectHandle food )
	<<
		-- TODO: anything w/ the particulars of food here
		state.AI.ints.hunger = 0 
	>>
	
	respond RequestMeleeAttackDamage()
	<<
		if state.entityData.melee_attack then
			-- BTW: 4th returned variable is for "on hit effect"
			local melee_sound = ""
			if state.entityData.sounds and state.entityData.sounds.attack then
				state.entityData.sounds = state.entityData.sounds.attack
			end
			return "MeleeAttackDamageResponse",
				state.entityData.melee_attack.damageType,
				state.entityData.melee_attack.damageAmount,
				state.entityData.melee_attack.onHitEffect,
				melee_sound
		else
			-- default values
			return "MeleeAttackDamageResponse",
				"blunt",
				3,
				"",
				"Animal Medium Attack"
		end
	>>
	
	receive Vocalize(string vocalization)
	<<
		--printl("animal", state.AI.name .. " received Vocalize: " .. vocalization )
		if state.entityData.sounds then
			if state.entityData.sounds[vocalization] then
				sound = state.entityData.sounds[vocalization]
				
				send("rendInteractiveObjectClassHandler",
					"odinRendererPlaySFXOnInteractive",
					state.renderHandle,
					sound )
			end
		end
	>>
	
	receive Nightfall()
	<<
		if SELF.tags["dead"] then return end
		-- It's transitioning to nighttime! Do stuff you'd do at night.
	>>
	
	receive setHungerTo( int hunger)
	<<
		state.AI.ints.hunger = hunger
	>>
	
	receive makeHostile()
	<<
		if SELF.tags.dead then return end
		
		send("gameSpatialDictionary",
			"gameObjectAddBit",
			SELF,
			8)
	>>
	
	receive NestCleared()
	<<
		SELF.tags.has_nest = nil
	>>
	
	receive AssignmentCancelledMessage( gameSimAssignmentHandle assignment )
	<<
		printl("ai_agent", state.AI.name .. " received AssignmentCancelledMessage")
		send("rendInteractiveObjectClassHandler",
			"odinRendererRemoveInteraction",
			state.renderHandle,
			"Cancel orders for animals")

		state.assignment = nil
		send(SELF,"resetInteractions")
	>>

	receive JobCancelledMessage(gameSimJobInstanceHandle job)
	<<
		printl("ai_agent", state.AI.name .. " received JobCancelledMessage")
		state.assignment = nil
		send(SELF,"resetInteractions")
	>>
	
	receive resetInteractions()
	<<
		send("rendInteractiveObjectClassHandler",
			"odinRendererClearInteractions",
			state.renderHandle)
			
		if not SELF.tags.dead and
			not SELF.tags.pet then
			
			send("rendInteractiveObjectClassHandler",
				"odinRendererAddInteractions",
				state.renderHandle,
				"Hunt " .. state.entityName ,
				"Hunt Animals",
				"Hunt Animals", --"Hunt Animals",
				"Hunt Animals", --"Hunt Animal (player order)",
				"hunting", -- icon
				"workshop", -- filter
				"Musket Fire", -- sound
				true,false)	
			
		elseif SELF.tags.dead and
			SELF.tags.meat_source then
			
			send("rendInteractiveObjectClassHandler",
				"odinRendererAddInteractions",
				state.renderHandle,
				"Butcher",
				"Butcher Corpse (player order)",
				"", --"Butcher",
				"", --"Butcher Corpse (player order)",
				"corpse",
				"",
				"Slice Flesh",
				true,true)
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
		if messagereceived == "Hunt Animals" and
			not state.assignment then
			
			send("gameBlackboard","gameObjectRemoveTargetingJobs",	SELF,nil)
			
			send("rendCommandManager",
				"odinRendererCreateParticleSystemMessage",
				"Small Beacon",
				state.AI.position.x,
				state.AI.position.y)
			
			if not assignment then
				assignment = query("gameBlackboard",
					"gameObjectNewAssignmentMessage",
					SELF,
					"Hunt Animals",
					"",
					"")[1]
					--"hunting",
					--"hunting")[1]
			end
			
			send("gameBlackboard",
				"gameObjectNewJobToAssignment",
				assignment,
				SELF,
				"Hunt Animals",
				"animal",
				true )
			
			send("rendOdinCharacterClassHandler",
				"odinRendererCharacterExpression",
				state.renderHandle,
				"dialogue",
				"hunting",
				true)
			
			send("rendInteractiveObjectClassHandler",
				"odinRendererClearInteractions",
				state.renderHandle)
			
			state.assignment = assignment
			setCancelInteraction = true
			
		elseif messagereceived == "Butcher Corpse (player order)" and
			not state.assignment then
			send("gameBlackboard","gameObjectRemoveTargetingJobs",	SELF,nil)

			send("rendCommandManager",
				"odinRendererCreateParticleSystemMessage",
				"Small Beacon",
				state.AI.position.x,
				state.AI.position.y)
			
			if not assignment then
				assignment = query("gameBlackboard",
					"gameObjectNewAssignmentMessage",
					SELF,
					"Butcher Animal",
					"",
					"")[1]
			end
			
			send("gameBlackboard",
				"gameObjectNewJobToAssignment",
				assignment,
				SELF,
				"Butcher Corpse (player order)",
				"corpse",
				true )
			
			send("rendOdinCharacterClassHandler",
				"odinRendererCharacterExpression",
				state.renderHandle,
				"dialogue",
				"butcher",
				true)
			
			send("rendInteractiveObjectClassHandler",
				"odinRendererClearInteractions",
				state.renderHandle)
			
			state.assignment = assignment
			setCancelInteraction = true
			
		elseif messagereceived == "Cancel orders for animals" and
			state.assignment then
			send("gameBlackboard",
				"gameObjectRemoveTargetingJobs",
				SELF,
				nil)
			
			state.assignment = nil
			send(SELF,"resetInteractions")
		end
		
		if setCancelInteraction then
			send("rendInteractiveObjectClassHandler",
                    "odinRendererAddInteractions",
				state.renderHandle,
                         "Cancel orders for animal",
                         "Cancel orders for animals",
                         "", --"Cancel orders for animals",
                         "", --"Cancel orders for animals",
						"",
						"",
						"",
						false,true)
		end
	>>
>>