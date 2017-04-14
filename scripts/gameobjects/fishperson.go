gameobject "fishperson" inherit "ai_agent"
<<
	local
	<<
          function fishperson_doOneSecondUpdate()
			
               -- normalize
               if state.AI.ints["morale"] < 0 then
                    state.AI.ints["morale"] = 0
               elseif state.AI.ints["morale"] > 100 then
                    state.AI.ints["morale"] = 100
               else
				-- thx Samut
                    local isNight = query("gameSession", "getSessionBool", "isNight")[1]
				if isNight then
					state.AI.ints["morale"] = state.AI.ints["morale"] + rand(1,2)
				else
					state.AI.ints["morale"] = state.AI.ints["morale"] + rand(0,1)
				end
               end
               
               if SELF.tags["fleeing"] and state.AI.ints["morale"] >= 75 then 
                    SELF.tags["fleeing"] = false
               end
          end
	>>

	state
	<<
		gameAIAttributes AI
          table entityData
          table traits
          string entityName
		string animSet
		int renderHandle
		bool asleep
		gameObjectHandle group
	>>

	receive Create( stringstringMapHandle init )
	<<
		local entityName = init["legacyString"]
          if not entityName then
               printl("ai_agent", " fishperson: improper legacystring / entityname given!")
               return
		elseif entityName == "Fishperson" then
			entityName = "fishperson"
          end
		
          state.entityName = entityName
		
		local entityData = EntityDB[ state.entityName ]
		if not entityData then
               printl("ai_agent", "fishperson: looking for entityName: " .. tostring(state.entityName) .. " but found nothing!" )
			return
		end
          state.entityData = entityData

		local fishInfo = EntityDB["FishInfo"]
		local function getPhoneme()
			return fishInfo.fishPhonemes[ rand(1,#fishInfo.fishPhonemes) ]
		end
		
		local name = ""
		for i=1, rand(2,3), 1 do
			if (rand(1,3) == 1) and (i > 1) then
				name = name .. "'" .. getPhoneme()
			else
				name = name .. getPhoneme()
			end
		end
		name = name:gsub("^%l", string.upper)
		
		state.AI.name = name
		
		state.AI.strs.firstName = name
		state.AI.strs.lastName = "the Fishperson"
		
		if init.firstName then
			state.AI.strs.firstName = firstName
		end
		if init.lastName then
			state.AI.strs.lastName = lastName
		end
		
		state.AI.strs["citizenClass"] = state.entityName

          -- START ai_damage required stats
          state.AI.ints["healthMax"] = 8
          state.AI.ints["healthTimer"] = 6 -- in seconds, per 1 point 
          state.AI.ints["fire_timer"] = 10
          state.AI.ints["health"] = state.AI.ints["healthMax"]
          state.AI.ints["numAfflictions"] = 0
          -- END ai_damage required stats
          
          state.AI.ints["object_attack_counter"] = 0
		
          state.AI.ints.aggression = rand(-2,3)
          state.AI.ints["morale"] = 100
		
		local humanstats = EntityDB["HumanStats"]
		local worldstats = EntityDB["WorldStats"]
		
          state.AI.ints["corpse_timer"] = humanstats["corpseRotTimeDays"] * worldstats["dayNightCycleSeconds"] * 10 -- in gameticks
		state.AI.ints.corpse_vermin_spawn_time_start = div(state.AI.ints.corpse_timer,2)
          state.AI.walkTicks = 3
          state.AI.ints["subGridWalkTicks"] = state.AI.walkTicks
          setposition(0,0)

          state.animSet = state.entityData.animationSetM 

		-- copy tag list over
		SELF.tags = state.entityData.tags
          
          -- no particular order to whether fishpeople have clothes or not YET, so:
          if rand(1,2) == 1 then
               SELF.tags["fishperson_clothed"] = true
          end
		
          local colours = { [1]="blue", [2]="grey", [3]="green" }
          local colour = colours[ rand(1,3) ]
		if init.colour then
			colour = init.colour
		end
          local bodycolour = colour

          if SELF.tags["fishperson_clothed"] then
               bodycolour = colour .. "_clothed"
          end
          
		local hat = ""
		local head = state.entityData.headsM[colour]
		
		if init.leader then
			-- hat or helmet? On a true leader can wear one.
			if rand(1,3) == 3 then
				head =state.entityData.helmets[ rand(1,#state.entityData.helmets) ]
			else
				hat = state.entityData.hats[ rand(1,#state.entityData.hats) ]
			end
		end
		
          send( "rendOdinCharacterClassHandler",
                    "odinRendererCreateCitizen",
                    SELF,
                    state.entityData.modelsM[bodycolour],
                    head, -- head
                    "", -- hair
				hat, -- hat
                    state.entityData.animationSetM,
                    0, 0 )
	
          -- set up simple personality until someone has a better idea
          local personalities = { "sanguine", "choleric", "melancholic", "phlegmatic" }
          SELF.tags[ personalities[rand(1,4)] ] = true
     
          state.AI.walkTicks = entityData.walkTicks
          send( "rendOdinCharacterClassHandler", "odinRendererSetCharacterWalkTicks", 
               SELF.id, state.AI.walkTicks )
		
		-- NOT UNTIL FOUL DAGON WALKS THE EARTH AGAIN SHALL THE FISHPEOPLE DEIGN TO JOIN US
          -- setting this postive seems to make them not like humans, actually.
		send( "gameAIPreferences", "changeAITagPreference", SELF, "human", 10)

		local occupancyMap = 
			".-.\\".. 
			"-C-\\"..
			".-.\\"
		local occupancyMapRotate45 = 
			".-.\\".. 
			"-C-\\"..
			".-.\\"
			
		send( "gameSpatialDictionary",
			"registerSpatialMapString",
			SELF,
			occupancyMap,
			occupancyMapRotate45,
			true )
			
          state.AI.ints["emoteTimer"] = rand(5,15)
          state.timer = 0
		
		local hostilecheck = query("gameSession", "getSessionBool", "fishpeoplePolicyHostile")[1]
		if hostilecheck then
			send(SELF, "setAggression", rand(1,3) )		
			send(SELF,"makeHostile")
			send("gameSpatialDictionary", "gameObjectAddBit", SELF, 10)
			send(SELF,"addTag","hostile_horror")
		else
			send(SELF, "setAggression", rand(-2,3) )
			send(SELF,"makeNeutral")
		end
		
		-- accounting.
		send("gameSession","incSessionInt","fishpeoplePopulation",1)
		
		-- If I am a hostile fishperson, I am hostile to all other players.
		
		-- Kill 'em all
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 0)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 1)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 2)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 3)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 4)
		

		-- FIXME nations or whatever

          ready()
		wake()
         
	>>

	receive GameObjectPlace( int x, int y ) 
	<<
		-- equip random melee & ranged weapon.
		if not state.AI.bools.first_placement_done then 
			local entity_data = EntityDB[ state.entityName ]
			
			send(SELF,"setWeapon","melee",entity_data.melee_weapons[ rand(1,#entity_data.melee_weapons) ] )
			send(SELF,"setWeapon","ranged",entity_data.ranged_weapons[ rand(1,#entity_data.ranged_weapons) ])
			
			state.AI.bools.first_placement_done = true
		end
		send(SELF,"resetInteractions")
	>>

	receive gameFogOfWarExplored(int x, int y )
	<<
		wake()
		--state.asleep = false
	>>

	receive makeHostile()
	<<
		printl("ai_agent", state.AI.name .. " : making this fisheprson hostile.")
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle,
			"name",
			state.AI.name)
		
		-- make me fishperson faction
		send("gameSpatialDictionary","gameObjectAddBit",SELF,10) -- Fishpeople
		
		-- and make me hostile to all players/the empire / everything.
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 0)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 1)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 2)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 3)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 4)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 5)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 6)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 7)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 8) -- Carnivores

		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 11) -- Obeliskians
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 14) -- Bandits
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 15) -- Frontier Justice Targets
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 16) -- Cultist Murder Targets

		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterCustomTooltipMessage",
			SELF.id,
			"ui\\tooltips\\fishpersonHostileTooltip.xml")
	>>
	
	receive makeFriendly()
	<<
		-- TODO: only give proper name if translation or friendly relations or something
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle,
			"name",
			state.AI.name)

		send("gameSpatialDictionary","gameObjectRemoveHostileBit",SELF,0)
		send("gameSpatialDictionary","gameObjectRemoveHostileBit",SELF,1)
		send("gameSpatialDictionary","gameObjectRemoveHostileBit",SELF,2)
		send("gameSpatialDictionary","gameObjectRemoveHostileBit",SELF,3)

		send("gameSpatialDictionary","gameObjectRemoveBit",SELF,10)

		SELF.tags["hostile_agent"] = nil
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterCustomTooltipMessage",
			SELF.id,
			"ui\\tooltips\\fishpersonFriendlyTooltip.xml")
	>>
	
	receive makeNeutral()
	<<
		-- TODO: only give proper name if translation or friendly relations or something
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle,
			"name",
			state.AI.name)
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterCustomTooltipMessage",
			SELF.id,
			"ui\\tooltips\\fishpersonTooltip.xml")
	>>

	receive Update()
	<<
		tooltip_refresh_from_save()
		
		if state.AI.thinkLocked then
               return
          end

          if state.AI.ints.updateTimer % 10 == 0 then
               fishperson_doOneSecondUpdate()
          end
    
		if state.AI.bools["dead"] then
			if not SELF.tags["buried"] then
				send(SELF,"corpseUpdate")
			else
				disable_buried_corpses() -- ai_agent.go function
			end	
			return
		end

          -- get a job, fishy!
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
				if VALUE_STORE[ "VerboseFSM" ] then
					if VALUE_STORE["showFSMDebugConsole"] then printl("FSM", "Fishperson Update #" .. tostring(SELF) .. ": received job " .. state.AI.curJobInstance.name) end

				end
			end
          else
               -- interrupt only at 1s intervals because fishpeople getting stuff right isn't as important as humans
               local oneSecond = false
               if state.AI.ints.updateTimer % 10 == 0 then
                    oneSecond = true
               end
               
               if oneSecond then 
                    local results = query( "gameBlackboard", "gameAgentTestForInterruptsMessage", state.AI, SELF )
                    if results.name == "gameAgentAssignedJobMessage" then
                         results[1].assignedCitizen = SELF
                         if state.AI.curJobInstance then
                              if VALUE_STORE["showFSMDebugConsole"] then printl("FSM", "FSM: Attempting to abort job " .. state.AI.curJobInstance.displayName .. " due to an interrupt!" ) end

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
	>>

	receive Vocalize(string vocalization)
	<<
		--printl("fishperson", state.AI.name .. " received Vocalize: " .. vocalization )
		if state.entityData.sounds then
			if state.entityData.sounds[vocalization] then
				sound = state.entityData.sounds[vocalization]
				
				send("rendInteractiveObjectClassHandler",
					"odinRendererPlaySFXOnInteractive",
					state.renderHandle,
					sound )
				
			elseif vocalization == "Converse" then
				sound = state.entityData.sounds["idle"]
				send("rendInteractiveObjectClassHandler",
					"odinRendererPlaySFXOnInteractive",
					state.renderHandle,
					sound )
				
			end
		end
	>>

	receive deathBy( gameObjectHandle damagingObject, string damageType )
	<<
		send(SELF,"Vocalize","die") 
		
		-- TODO flesh out handling of damagingObject and damageType into interesting descriptions.
		SELF.tags["horror"] = nil
		SELF.tags.horror_corpse = true
		SELF.tags["hostile_horror"] = nil
		SELF.tags.meat_source = true
		
		-- explode into meats if blown up
		if damageType == "explosion" or damageType == "shrapnel" then
			meat_splosion()
		end

		send("gameSpatialDictionary","gameObjectRemoveBit",SELF,10)
		
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
	
		send("gameSession","incSessionInt","fishpeoplePopulation",-1)
		send("gameSession","incSessionInt","fishpeopleDeaths", 1)
		
		local fishpeople_killed = query("gameSession","getSessionInt","fishpeopleDeaths")[1]
		if not query("gameSession","getSessionBool","killedManyFishpeople")[1] and fishpeople_killed >= 100 then
			send("gameSession", "setSessionBool", "killedManyFishpeople", true)
		end
		send("gameSession", "incSteamStat", "stat_fishpeople_killed", 1)
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterCustomTooltipMessage",
			SELF.id,
			"ui\\tooltips\\fishpersonDeadTooltip.xml")
			
		send(SELF,"resetInteractions")
		
		local civ = query("gameSpatialDictionary", "gridGetCivilization", state.AI.position )[1]
		if civ == 0 then
			if query("gameSession","getSessionBool","horror_policy_dump")[1] then
				send(SELF,
					"HandleInteractiveMessage",
					"Dump Corpse (player order)",
					nil)
				
			elseif query("gameSession","getSessionBool","horror_policy_harvest")[1] then
				-- it'll be done automatically.
				--[[send(SELF,
					"HandleInteractiveMessage",
					"Butcher Fishperson Corpse (player order)",
					nil)
				]]
			elseif query("gameSession","getSessionBool","horror_policy_study")[1] then
				send(SELF,
					"HandleInteractiveMessage",
					"Study Horror (fishperson)",
					nil)
			end
		end
	>>

	receive HarvestMessageNoProduct( gameObjectHandle harvester, gameSimJobInstanceHandle ji )
	<<
		--send("gameBlackboard", "gameObjectRemoveTargetingJobs", SELF, ji)
		
		-- Put some music in the dynamics.
		incMusic(3,10)
		--send(SELF, "ForceDropTools")
		
		-- TODO: set up fishperson skeleton swap (if we have a model for it?)
		send("rendOdinCharacterClassHandler",
			"odinRendererDeleteCharacterMessage",
			state.renderHandle)
		
		send("gameSpatialDictionary", "gridRemoveObject", SELF)
		destroyfromjob(SELF,ji)
	>>

	receive HarvestMessage( gameObjectHandle harvester, gameSimJobInstanceHandle ji )
	<<
		local numSteaks = 2 -- so gross.
		for s=1, numSteaks do
			local results = query( 	"scriptManager",
								"scriptCreateGameObjectRequest",
								"item",
								{legacyString = "raw_fishperson_steak"} )
			
			local handle = results[1]
			
			if( handle == nil ) then 
				--printl("Creation failed")
				return "abort"
			else 
			 	local ownerTags = query(harvester,"getTags")[1]	
				if ownerTags.citizen then
					send( handle,"ClaimItem")
				end

				send(handle,
					"GameObjectPlace",
					state.AI.position.x,
					state.AI.position.y  )
			end
		end
		
		-- TODO: set up fishperson skeleton swap (if we have a model for it?)
		send("rendOdinCharacterClassHandler",
			"odinRendererDeleteCharacterMessage",
			state.renderHandle)
		
		send("gameSpatialDictionary",
			"gridRemoveObject",
			SELF)
		
		destroyfromjob(SELF,ji)
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
     
     receive damageMessage( gameObjectHandle damagingObject, string damageType, int damageAmount, string onhit_effect )
     <<
		if SELF.tags.dead then
			return
		end
		
		if query("gameSession","getSessionBool","fishpeoplePolicyHostile")[1] == false then
			
			if SELF.tags["was_attacked"] ~= true then
				SELF.tags.was_attacked = true

				if state.AI.ints.aggression == 0 then
					state.AI.ints.aggression = state.AI.ints.aggression + rand(-1,1)
				end
				
				if damagingObject then 
					local damagerTags = query(damagingObject, "getTags")[1]

					if damagerTags.citizen then
						-- do we switch to attack more or not?
						if state.AI.ints.aggression > 1 then
							SELF.tags.hostile_agent = true
							SELF.tags.hostile_horror = true
						else
							send(SELF,"AICancelJob", "intimidation")
							send( "gameBlackboard",
								"gameCitizenJobToMailboxMessage",
								SELF,
								damagingObject,
								"Flee Due To Intimidation",
								"enemy")
						end
						
						-- send exclamation to fellows.
						--[[local results = query("gameSpatialDictionary",
										  "allObjectsInRadiusWithTagRequest",
										  state.AI.position,
										  10,
										  "fishperson",
										  true)--]]
						
						local results = query("gameSpatialDictionary",
									 "allObjectsInRadiusRequest",
									 state.AI.position,
									 10,
									 true)
						
						--[[if results and results[1] then
							send(results[1], "detectCorpse", SELF)
						end--]]
						
						if results then
							send(results[1],
								"hearExclamation",
								"fishperson_attacked",
								SELF,
								damagingObject)
						end
					end
				end
			end
			
			--if damagerTags.citizen and damagerTags.combat_target_for_enemy then
			--[[if state.AI.ints.aggression > 0 then
				-- attacked by legitimate colonist; hurt them!
				send( "gameBlackboard",
					"gameCitizenJobToMailboxMessage",
					SELF,
					damagingObject,
					"Fishperson Retaliate Against Attacker",
					"enemy")
				
			else
				send( "gameBlackboard",
					"gameCitizenJobToMailboxMessage",
					SELF,
					damagingObject,
					"Flee Due To Intimidation",
					"enemy")
			end--]]
			--end
		end
		-- decrease morale because getting hurt is scary
		state.AI.ints["morale"] = state.AI.ints["morale"] - damageAmount * 3
		if state.AI.ints["morale"] <= 50 and not SELF.tags["fleeing"] and not SELF.tags["rampage"] then
			FSM.abort(state,"Morale broken." )
			SELF.tags["fleeing"] = true
			SELF.tags["rampage"] = false
			SELF.tags["raider"] = false
		end
     >>
     
	receive spawnGibs()
	<<  
		if state.entityData["gibs"] then
			for s=1, rand( state.entityData["gibs"].min , state.entityData["gibs"].max ) do
                local gibName = state.entityData["gibs"].name
				results = query( "scriptManager", "scriptCreateGameObjectRequest", "clearable", { legacyString = gibName } )
				handle = results[1]
				
                    if not handle then 
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
		
		if not state.AI.bools["rotted"] then
			
			state.AI.ints["corpse_timer"] = state.AI.ints["corpse_timer"] - 1
		
			-- broadcast that there's a rotting corpse over here.
			-- oh, and make a corpse return job if we haven't
			if state.AI.ints["corpse_timer"] % 100 == 0 then
                    -- timer in gameticks, trigger once per 10s
				local results = query("gameSpatialDictionary",
							 "allObjectsInRadiusRequest",
							 state.AI.position,
							 10,
							 true)
				
				if results and results[1] then
					send(results[1], "detectCorpse", SELF)
				end
				
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
				
				-- no.
				--send(SELF,"makeDumpCorpseJob")
			end

			if state.AI.ints["corpse_timer"] <= 0 then
				-- here's your skeleton model swap
				state.AI.bools["rotted"] = true
				SELF.tags["meat_source"] = nil
				
				-- TODO : detach particles from joint. No working examples of this in code.
				-- FIXME : the animation explodes here for some reason. Removing for now.  (It's because we're not checking for biped or dress -DJ)
                    
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
                    
				send("rendOdinCharacterClassHandler",
					"odinRendererDeleteCharacterMessage",
					state.renderHandle)
				
				send("gameSpatialDictionary", "gridRemoveObject", SELF)
				destroy(SELF)
			end
		end
	>>
     
	receive hearExclamation( string name, gameObjectHandle exclaimer, gameObjectHandle subject )
	<<
		if SELF.tags["dead"] then return end
		
		if name == "explosion" then
			local explosionName =  query(exclaimer, "getName")[1]
			
			if (explosionName ~= "Small Explosion") and (explosionName ~= "Urchin Grenade Explosion") and (explosionName ~= "Blunderbuss Area Attack") then
				state.AI.ints["morale"] = state.AI.ints["morale"] - 5
				--send(SELF, "attemptEmote", "explosion",4,false)
				
				if state.AI.ints["morale"] <= 50 and not SELF.tags["fleeing"] and not SELF.tags["rampage"] then
					--FSM.abort(state,"Morale broken.")
					SELF.tags["fleeing"] = true
					SELF.tags["rampage"] = false
					SELF.tags["raider"] = false
				end
			end
		elseif name == "detectCorpse" then
			state.AI.ints["morale"] = state.AI.ints["morale"] - 3
			if state.AI.ints["morale"] <= 50 and not SELF.tags["fleeing"] and not SELF.tags["rampage"] then
				--FSM.abort(state,"Morale broken.")
				SELF.tags["fleeing"] = true
				SELF.tags["rampage"] = false
				SELF.tags["raider"] = false
			end
		elseif name == "detectCombat" then
			wake()
			-- decrease morale because combat is scary
			state.AI.ints["morale"] = state.AI.ints["morale"] - 3
			if VALUE_STORE["showCombatDebugConsole"] then printl("combat", state.AI.name .. " the fishperson witnessed combat, morale now: " .. state.AI.ints["morale"] ) end
			if state.AI.ints["morale"] <= 50 and not SELF.tags["fleeing"] and not SELF.tags["rampage"] then
				--FSM.abort(state,"Morale broken.")
				SELF.tags["fleeing"] = true
				SELF.tags["rampage"] = false
				SELF.tags["raider"] = false
			end
		elseif name == "fishperson_attacked" then
			
			-- a fellow fishperson is attacked!
			-- either become more passive or more aggressive.
			
			if state.AI.ints.aggression >= 2 then
				-- fight!
				state.AI.ints.aggression = state.AI.ints.aggression + 1
				SELF.tags.hostile_agent = true
				SELF.tags.hostile_horror = true
				send(SELF, "attemptEmote", "fishperson_bloody", 4,true )
			elseif state.AI.ints.aggression == 1 then
				-- anger emote.
				state.AI.ints.aggression = state.AI.ints.aggression + 1
				send(SELF, "attemptEmote", "fishperson_angry", 4,true )
			elseif state.AI.ints.aggression == 0 then
				-- fear emote
				state.AI.ints.aggression = state.AI.ints.aggression - 1
				send(SELF, "attemptEmote", "fishperson_scared", 4,true )
			end
			
			if state.AI.ints.aggression < 0 then
				-- run!
				send(SELF,"AICancelJob", "intimidation")
				send( "gameBlackboard",
					"gameCitizenJobToMailboxMessage",
					SELF,
					damagingObject,
					"Flee Due To Intimidation",
					"enemy")
			end
		elseif name == "fishpeople_eggs_plundered" then
			state.AI.ints.aggression = state.AI.ints.aggression + 1
			--send(SELF, "attemptEmote", "fishperson_angry", 4,true )
			if SELF.tags.sanguine then
				-- go hostile
				if state.group then
					send(state.group,
						"removeMember",
						SELF,
						"Taking solo revengne upon humanity.")
				end
				
				send(SELF,"makeHostile")
				send("rendOdinCharacterClassHandler",
					"odinRendererCharacterExpression",
					state.renderHandle,
					"thought",
					"fishperson_bloody",
					true)
				
				send("gameBlackboard",
					"gameCitizenJobToMailboxMessage",
					member,
					nil,
					"Charge Toward Civilization (fishperson)",
					"")
				
			elseif SELF.tags.choleric then
				-- form raiding party
				send(SELF,"attemptEmote","fishpeople_attack",3,true)
				if state.group then
					send("rendCommandManager",
						"odinRendererTickerMessage",
						"Fishpeople discovered that their eggs were plundered by humans! Now they're coming after you to enact bloody vengeance.",
						"fishperson_bloody",
						"ui\\thoughtIcons.xml")
					
					send(state.group,"pushMission","attack",nil,30)
				end
				SELF.tags.rallying_raiders = nil
				
			elseif SELF.tags.melancholic then
				-- cry
				send("rendOdinCharacterClassHandler",
					"odinRendererCharacterExpression",
					state.renderHandle,
					"thought",
					"fishperson_scared",
					true)
				
			end
		end
	>>
     
     receive resetEmoteTimer()
	<<
		state.AI.ints["emoteTimer"] = 0
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
			--destroyfromjob(SELF, ji)
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
	
	receive setAggression(int aggro)
	<<
		state.AI.ints.aggression = aggro
	>>
	
	receive incAggression(int aggro)
	<<
		state.AI.ints.aggression = state.AI.ints.aggression + aggro
	>>
	
	receive resetAggression()
	<<
		-- TODO: decide to use this or delete it.
		-- factors: group's mission, overall fishpeople relations, fishperson traits.
	>>

	receive Nightfall()
	<<
		-- It's transitioning to nighttime! Do stuff you'd do at night.
		if SELF.tags["dead"] then return end
		
		if state.AI then
			if state.AI.ints.hunger then 
				-- NOTE: not used yet by fishpeople, but in case we care ... 
				state.AI.ints.hunger = state.AI.ints.hunger + 1
				state.AI.ints.tiredness = state.AI.ints.tiredness + 1
			end
		end
	>>
	
	receive AssignmentCancelledMessage( gameSimAssignmentHandle assignment )
	<<
		printl("ai_agent", state.AI.name .. " received AssignmentCancelledMessage")
		send("rendInteractiveObjectClassHandler",
			"odinRendererRemoveInteraction",
			state.renderHandle,
			"Cancel corpse orders")

		state.assignment = nil
		send(SELF,"resetInteractions")
	>>

	receive JobCancelledMessage(gameSimJobInstanceHandle job)
	<<
		printl("ai_agent", state.AI.name .. " received JobCancelledMessage")
		state.assignment = nil
		send(SELF,"resetInteractions")
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
		
		if messagereceived == "Attack Fishpeople"  and
			not state.assignment and
			not SELF.tags.dead then
			
			-- TODO: do an event for this flip to hostile.
			if state.group then
				send(state.group,"makeHostile")
			else
				send(SELF,"makeHostile")
			end
			
			send("gameSession", "incSessionInt", "fishpeopleRelations", -5)
			
			send("gameSession","setSessionBool","fishpeoplePolicyHostile",true)
			send("gameSession","setSessionBool","fishpeoplePolicyDenial",false)
			send("gameSession","setSessionBool","fishpeoplePolicyFriendly",false)
			
			send("rendCommandManager",
				"odinRendererTickerMessage",
				"You've ordered your troops to attack the Fishpeople, starting with " .. state.AI.name .. "!",
				"fishperson_angry",
				"ui\\thoughtIcons.xml")
			
			-- Can't undo this one!
			setCancelInteraction = false
			state.assignment = assignment
				
		elseif messagereceived == "Hassle Fishpeople" and
			not state.assignment and
			not SELF.tags.dead then
			
			-- jobs will take it from here.
			SELF.tags.bad_fishperson = true
			
			setCancelInteraction = false -- true
			
		--[[elseif message == "Share Food With Fishpeople" and
			not SELF.tags.invited_to_dinner then
			
			if state.group then
				send(state.group,"shareFoodWith")
			else
				SELF.tags.invited_to_dinner = true
			end--]]
			
		elseif messagereceived == "Study Horror (fishperson)"  and
			not state.assignment and
			SELF.tags.dead then
			
			send("gameBlackboard",
				"gameObjectRemoveTargetingJobs",
				SELF,
				nil)
			
			send("rendCommandManager",
				"odinRendererCreateParticleSystemMessage",
				"Small Beacon",
				state.AI.position.x,
				state.AI.position.y)
			
			if not assignment then
				assignment = query("gameBlackboard",
					"gameObjectNewAssignmentMessage",
					SELF,
					"Study Fishperson Corpse",
					"",
					"")[1]
			end
				
			send( "gameBlackboard",
				"gameObjectNewJobToAssignment",
				assignment,
				SELF,
				"Study Horror (fishperson)",
				"corpse",
				true )
			
			-- thoughtIcons.xml only.
			send("rendOdinCharacterClassHandler",
				"odinRendererCharacterExpression",
				state.renderHandle,
				"dialogue",
				"fishperson_dead",
				true)
			
			send("rendInteractiveObjectClassHandler",
				"odinRendererClearInteractions",
				state.renderHandle)
			
			setCancelInteraction = false -- true
			state.assignment = assignment
			
		elseif messagereceived == "Butcher Fishperson Corpse (player order)"  and
			not state.assignment and
			SELF.tags.dead then
			
			send("gameBlackboard",
				"gameObjectRemoveTargetingJobs",
				SELF,
				nil)
			
			send("rendCommandManager",
				"odinRendererCreateParticleSystemMessage",
				"Small Beacon",
				state.AI.position.x,
				state.AI.position.y)
			
			if not assignment then
				assignment = query("gameBlackboard",
					"gameObjectNewAssignmentMessage",
					SELF,
					"Butcher Fishperson Corpse",
					"",
					"")[1]
			end
				
			send( "gameBlackboard",
				"gameObjectNewJobToAssignment",
				assignment,
				SELF,
				"Butcher Fishperson Corpse (player order)",
				"corpse",
				true )
			
			-- thoughtIcons.xml only.
			send("rendOdinCharacterClassHandler",
				"odinRendererCharacterExpression",
				state.renderHandle,
				"dialogue",
				"butcher",
				true)
			
			send("rendInteractiveObjectClassHandler",
				"odinRendererClearInteractions",
				state.renderHandle)
			
			setCancelInteraction = false -- true
			state.assignment = assignment
	
		elseif messagereceived == "Bury Corpse (player order)" and
			not state.assignment and
			SELF.tags.dead then
			
			-- TODO: fire assignment for controversy over burying fishpeople
			
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
			not state.assignment and
			SELF.tags.dead then
			
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
				true)
			
			send("rendOdinCharacterClassHandler",
				"odinRendererCharacterExpression",
				state.renderHandle,
				"dialogue",
				"jobhand",
				true)
			
			setCancelInteraction = true
			state.assignment = assignment
			
		--[[elseif messagereceived == "Cancel Fishperson orders" and
			state.assignment and
			SELF.tags.dead then
			
			send("gameBlackboard",
				"gameObjectRemoveTargetingJobs",
				SELF,
				nil)
			
			state.assignment = nil
			send(SELF,"resetInteractions")
			SELF.tags.bad_fishperson = false -- hacky.]]--
			
		elseif messagereceived == "Cancel corpse orders" and
			state.assignment and
			SELF.tags.dead then
			
			send("gameBlackboard",
				"gameObjectRemoveTargetingJobs",
				SELF,
				nil)
			
			send(SELF,"resetInteractions")
			state.assignment = nil
		end
		
		if setCancelInteraction then
			--if SELF.tags.dead then
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
							false,true)
			--[[else
				send("rendInteractiveObjectClassHandler",
						"odinRendererAddInteractions",
					state.renderHandle,
							 "Cancel orders regarding " .. state.AI.name,
							 "Cancel Fishperson orders",
							 "Cancel Fishperson orders", --"Cancel Fishperson orders",
							 "Cancel Fishperson orders", --"Cancel Fishperson orders",
							"",
							"",
							"",
							false,true)
			end]]--
		end
	>>
	
	receive resetInteractions()
	<<
		--printl("ai_agent", state.AI.name .. " receive resetInteractions")
		
		send("rendInteractiveObjectClassHandler",
			"odinRendererClearInteractions",
			state.renderHandle)

		--local denial =	query("gameSession","getSessionBool","fishpeoplePolicyDenial")[1]
		--local friendly = query("gameSession","getSessionBool","fishpeoplePolicyFriendly")[1]
		
		if SELF.tags.dead and not SELF.tags.buried then
			if not state.assignment then
				if SELF.tags.meat_source then 
					send("rendInteractiveObjectClassHandler",
						"odinRendererAddInteractions",
						state.renderHandle,
						"Butcher Fishperson",
						"Butcher Fishperson Corpse (player order)",
						"Butcher Fishpeople", --"Butcher Fishpeople",
						"Butcher Fishperson Corpse (player order)", --"Butcher Fishperson Corpse (player order)",
							"butcher_icon",
							"",
							"Flesh Crack",
							false,true)
					
					send("rendInteractiveObjectClassHandler",
						"odinRendererAddInteractions",
						state.renderHandle,
						"Have Naturalist Dissect Fishperson",
						"Study Horror (fishperson)",
						"Dissect Fishpeople",
						"Study Horror (fishperson)",
							"",
							"",
							"click",
							true,true)
				end
				
				send("rendInteractiveObjectClassHandler",
						"odinRendererAddInteractions",
					state.renderHandle,
							 "Dump Fishperson Corpse",
							 "Dump Corpse (player order)",
							 "Dump Corpses", --"Dump Corpses",
							 "Dump Corpse (player order)", --"Dump Corpse (player order)",
							"graveyard",
							"",
							"Dirt",
							false,true)
			else
				send("rendInteractiveObjectClassHandler",
						"odinRendererAddInteractions",
					state.renderHandle,
							 "Cancel orders for corpse of " .. state.AI.name,
							 "Cancel corpse orders",
							 "Cancel corpse orders", --"Cancel corpse orders",
							 "Cancel corpse orders", -- "Cancel corpse orders",
							"graveyard",
							"",
							"Dirt",
							false,true)				
			end
		elseif not SELF.tags.dead then
			
			-- not dead.
			local hostile = query("gameSession","getSessionBool","fishpeoplePolicyHostile")[1]
			
			if not hostile then
				send("rendInteractiveObjectClassHandler",
					"odinRendererAddInteractions",
					state.renderHandle,
					"Attack Fishpeople",
					"Attack Fishpeople",
					"", --"Attack Fishpeople",
					"", --"Attack Fishpeople)",
					"",
					"",
					"click",
					true,true)
				
				send("rendInteractiveObjectClassHandler",
					"odinRendererAddInteractions",
					state.renderHandle,
					"Hassle Fishpeople",
					"Hassle Fishpeople",
					"", --"Hassle Fishpeople",
					"", --"Hassle Fishpeople)",
					"",
					"",
					"click",
					true,true)
			end
		end
	>>
>>