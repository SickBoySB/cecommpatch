gameobject "quaggaroth" inherit "ai_agent"
<<
	local
	<<

	>>

	state
	<<
		gameObjectHandle group
		gameAIAttributes AI
		gameGridPosition spawnLocation
		gameGridPosition wanderDestination
		table entityData
		string entityName
		string animSet
		int renderHandle
		int timer
		bool asleep
	>>
	
	receive Create( stringstringMapHandle init )
	<<
		
		--[[local entityName = init["legacyString"]
		if not entityName then
			printl("ai_agent", "quaggaroth name not found: " .. tostring(entityName))
			return
		end]]
          state.entityName = "quaggaroth"
		
		local entityData = EntityDB[ state.entityName ]
		if not entityData then
			printl("ai_agent", "quaggaroth data not found")
			return
		end
          state.entityData = entityData
          
		state.renderHandle = SELF.id
          state.AI.strs["mood"] = "deconstructive"
		state.group = false
		
          state.AI.strs["firstName"] = "Quag'garoth "
		state.AI.strs["lastName"] = SELF.id
		
		state.AI.name = state.AI.strs.firstName .. " " .. state.AI.strs.lastName

		state.AI.strs["gender"] = "none"
		
		state.AI.strs["citizenClass"] = state.entityName
		state.model = state.entityData.model
		state.animSet = state.entityData.animationSet
		
		send( "rendOdinCharacterClassHandler",
			"odinRendererCreateCharacter", 
			SELF,
			state.model,
			state.animSet,
			0,
			0 )
		
		send("rendOdinCharacterClassHandler",
			"odinRendererFaceCharacter", 
			state.renderHandle, 
			state.AI.position.orientationX,
			state.AI.position.orientationY )

		SELF.tags = {}
		if state.entityData.job_classes then
			for k,v in pairs(state.entityData.job_classes) do
				SELF.tags[v] = true
			end
		end
		if state.entityData.tags then
			for k,v in pairs(state.entityData.tags) do
				SELF.tags[v] = true		
			end
		end

          -- START ai_damage required stats
		if state.entityData.health then
               state.AI.ints.healthMax = state.entityData.health.healthMax
               state.AI.ints.health = state.entityData.health.healthMax
               state.AI.ints.healthTimer = 0
               state.AI.ints.healthTimerMax = state.entityData.health.healthTimerSeconds
		else
			state.AI.ints.healthMax = 150
               state.AI.ints.health = 150
               state.AI.ints.healthTimer = 0
               state.AI.ints.healthTimerMax = 3
          end
		state.AI.ints["numAfflictions"] = 0
          state.AI.ints["fire_timer"] = 10
          -- END ai_damage required stats
          
          state.AI.ints["corpse_timer"] = -1 -- Quaggaroth doesn't rot.
		
		if state.entityData.walkTicks then
			state.AI.walkTicks = state.entityData.walkTicks
		else
	          state.AI.walkTicks = 3
		end
		
          state.AI.ints["subGridWalkTicks"] = state.AI.walkTicks
		
		setposition(0,0)
		
          send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterWalkTicks",
			SELF.id,
			state.AI.walkTicks )

		send("gameSpatialDictionary",
			"registerSpatialMapString",
			SELF,
			state.entityData.occupancyMap,
			state.entityData.occupancyMapRotate45,
			true )
		
          state.AI.ints["emoteTimer"] = rand(1,30)
          state.timer = 0
		  
		send("gameSpatialDictionary", "gameObjectAddBit", SELF, 11) -- obeliskian

		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 0) -- player 1
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 1) -- player 2
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 2) -- player 3
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 3) -- player 4
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 4) -- empire
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 5) -- republique
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 6) -- stahlmark
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 7) -- novorus
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 14) -- bandit
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle,
			"name",
			state.AI.name)
		
		--[[send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterCustomTooltipMessage",
			SELF.id,
			"ui\\tooltips\\steamknightTooltip.xml")]]
		
		
		
		state.time_until_leave = 1*60*10 -- 1-2 shifts?
		
		if init.exorcism_performed then
			state.time_until_leave = 60
		end
		
		if init.cult_leader_death then
			state.time_until_leave = 60*5
		end
		
		if init.shrine_destroyed then
			SELF.tags.vulnerable = true
		end
		
		if init.cult_power then
			state.AI.ints.healthMax = tonumber(init.cult_power) * 15
			if state.AI.ints.healthMax < 30 then
				state.AI.ints.healthMax = 30
			end
               state.AI.ints.health = state.AI.ints.healthMax
		end
			
		
          ready()
		
		if state.entityData.defaultMeleeAttack then
			send(SELF,"setWeapon","melee","quaggaroth")
		end
		
		if state.entityData.defaultRangedAttack then
			send(SELF,"setWeapon","ranged","quaggaroth")
		end
		
		--[[if init.dormant and init.dormant == "true" then
			send(SELF,"goDormant")
		elseif init.dormant and init.dormant == "false" then
			-- stay awake.
		else
			send(SELF,"goDormant")
		end]]
		
		send("rendOdinCharacterClassHandler",
			"odinRendererHideCharacterMessage",
			state.renderHandle,
			true)
		
		send("gameBlackboard",
			"gameCitizenJobToMailboxMessage",
			SELF,
			nil,
			"Burst From Ground (quaggaroth)",
			"")
		
	>>

	receive setGroup( gameObjectHandle group )
	<<
		state.group = group
		--[[if group ~= nil then
			SELF.tags.in_group = true
			SELF.tags.not_in_group = nil
		else
			SELF.tags.in_group = nil
			SELF.tags.not_in_group = true
			SELF.tags.active_mission = nil
			FSM.abort("disconnected from Master")
			send( "gameBlackboard",
				"gameCitizenJobToMailboxMessage",
				SELF,
				nil,
				"Disconnect From Assembly (Obeliskian)",
				"")
		end]]
	>>
	
	receive goDormant()
	<<
		--[[send( "gameBlackboard",
			"gameCitizenJobToMailboxMessage",
			SELF,
			nil,
			"Go Dormant (Obeliskian)",
			"")]]
	>>
	
	receive reactivate()
	<<
		--[[SELF.tags.dormant = nil
		send( "gameBlackboard",
			"gameCitizenJobToMailboxMessage",
			SELF,
			nil,
			"Reactivate (Obeliskian)",
			"")]]
	>>
	
	receive gameFogOfWarExplored(int x, int y )
	<<
		state.asleep = false
	>>

	receive GameObjectPlace( int x, int y ) 
	<<
		setposition(x,y)

		send("rendOdinCharacterClassHandler", 
				"odinRendererTeleportCharacterMessage", 
				state.renderHandle, 
				x,
				y )

		state.spawnLocation = gameGridPosition:new()
		state.spawnLocation.x = x
		state.spawnLocation.y = y
		
		--send("gameSession","incSessionInt","obeliskiansOnMap",1)
		send(SELF,"resetInteractions")
	>>
	
	receive SleepMessage()
	<<
		state.asleep = true
	>>

     receive Update()
     <<
		if not state or not state.AI or state.dead then
			return
		end
		
		send("gameSpatialDictionary",
			"gridExploreFogOfWar",
			state.AI.position.x,
			state.AI.position.y,
			6)
	
		if SELF.tags.dormant then
			return
		end
		
  		if state.AI.thinkLocked then
               return
          end
    
          state.timer = state.timer +1
          if state.timer % 10 == 0 then
               ai_agent_doOneSecondUpdate()
          end
		
		
		if state.time_until_leave > 0 then
			state.time_until_leave = state.time_until_leave -1
			
			if state.time_until_leave == 0 then
				SELF.tags.exit_map = true
			end
		end
		
		if not SELF.tags.dead and SELF.tags.trampling then
			if rand(1,5) == 1 then
				-- do AOE!
				-- DustPuffMassive
				 local handle = query("scriptManager",
                                        "scriptCreateGameObjectRequest",
                                        "explosion",
                                        { legacyString= "Quaggaroth Trample" } )[1]

				local positionResult = query("gameSpatialDictionary",
						   "nearbyEmptyGridSquare",
						   state.AI.position,
						   3)[1]
		
				if positionResult.onGrid then
					send(handle,
						"GameObjectPlace",
						positionResult.x,
						positionResult.y  )
				else
					send(handle,
						"GameObjectPlace",
						state.AI.position.x,
						state.AI.position.y  )
				end
					
			end
		end
    
		if state.AI.bools["dead"] then
               send(SELF, "corpseUpdate")
			SELF.tags.trampling = nil
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
			local results = query( "gameBlackboard", "gameAgentNeedsJobMessage", state.AI, SELF )

			if results.name == "gameAgentAssignedJobMessage" then 
				state.AI.curJobInstance = results[ 1 ]
				state.AI.curJobInstance.assignedCitizen = SELF
				send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "currentJob", state.AI.curJobInstance.displayName);	
				if VALUE_STORE[ "VerboseFSM" ] then
					if VALUE_STORE["showFSMDebugConsole"] then printl("FSM", "Citizen Update #" .. tostring(SELF) .. ": received job " .. state.AI.curJobInstance.name) end

				end
			end
          else
               -- interrupt only at 1s intervals because enemies getting stuff right isn't as important as humans
               local oneSecond = false
               if state.timer % 10 == 0 then
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
	>>
  
	respond getHomeLocation()
	<<
		-- used by a certain require_thing type hook in some jobs.
		-- look at bandits for example.
		
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
			"DustExplosion",
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
				elseif key ~= "curPickedUpItem" or not holdingChar or not state.AI.possessedObjects["curPickedUpCharacter"] then
					send(value, "DestroySelf", state.AI.curJobInstance )
				end
			end
		end
		
		-- destroy item if holding one.
		--send( state.AI.possessedObjects["curCarriedTool"],
		--	"DestroySelf",
		--	state.AI.curJobInstance )
		
		send("rendOdinCharacterClassHandler", "odinRendererDeleteCharacterMessage", state.renderHandle)
		send("gameSpatialDictionary", "gridRemoveObject", SELF)
		destroyfromjob(SELF,ji)
	>>

	receive deathBy( gameObjectHandle damagingObject, string damageType )
	<<
		SELF.tags["dead_horror"] = true
		SELF.tags.dead_quaggaroth = true
		SELF.tags["horror"] = nil
		SELF.tags["hostile_horror"] = nil
		
		send("gameSpatialDictionary","gameObjectRemoveBit",SELF,11)
		
		if state.AI.curJobInstance then
			FSM.abort( state, "Died.")
		end
		
		--[[send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterGeometry", 
			state.renderHandle,
			state.entityData.dormant_model,
			"",
			"none",
			state.animSet,
			"idle")]]

		-- animation time
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAnimationMessage",
			state.renderHandle,
			"death",
			false)
		
		local handle = query( "scriptManager",
					 "scriptCreateGameObjectRequest",
					 "explosion",
					 { legacyString = "Quaggaroth Explosion" } )[1]
		
		send(handle,
			"GameObjectPlace",
			state.AI.position.x,
			state.AI.position.y  )

		send(SELF,"resetInteractions")
		
		local handle = query( "scriptManager",
					 "scriptCreateGameObjectRequest",
					 "objectcluster",
					 { legacyString = "Obeliskian Cluster" } )[1]
		
		send(handle,
			"GameObjectPlace",
			state.AI.position.x + rand(-3,3),
			state.AI.position.y + rand(-3,3) )
		
		
		send("rendOdinCharacterClassHandler",
			"odinRendererDeleteCharacterMessage",
			state.renderHandle)
		
          send("gameSpatialDictionary",
			"gridRemoveObject",
			SELF)
		
          send("gameBlackboard",
			"gameObjectRemoveTargetingJobs",
			SELF,
			ji)
		
		destroyfromjob(SELF, nil)
	>>
	
	receive HarvestMessageNoProduct( gameObjectHandle harvester, gameSimJobInstanceHandle ji )
	<<
		send(SELF,"HarvestMessage",harvester, ji)
	>>
	
	receive HarvestMessage( gameObjectHandle harvester, gameSimJobInstanceHandle ji )
	<<
		send("rendOdinCharacterClassHandler",
			"odinRendererDeleteCharacterMessage",
			state.renderHandle)
		
		send("gameSpatialDictionary", "gridRemoveObject", SELF)
		destroyfromjob(SELF,ji)
	>>

     respond getIdleAnimQueryRequest()
	<<
		local animName = "idle" 
		local idleAnims = {"idle", "examine_below",}
		animName = idleAnims[rand(1,#idleAnims)]

		return "idleAnimQueryResponse", animName
	>>
     
	receive spawnGibs()
	<<
		-- TODO: pull this from EDB
		local handle = query("scriptManager",
						"scriptCreateGameObjectRequest",
						"objectcluster",
						{legacyString = "Obeliskian Gibs",})[1]
			
		send(handle,
			"GameObjectPlace",
			state.AI.position.x,
			state.AI.position.y)
	>>

     receive hearExclamation( string name, gameObjectHandle exclaimer, gameObjectHandle subject )
	<<
		if SELF.tags.dead then
			return
		end
		--[[
		if name == "explosion" or
			name == "mining" or
			name == "study_rock" then
			
			if query("gameSession", "getSessionInt", "dayCount")[1] > 4 then
				
				if rand(1,4) == 1 then
					if state.group then
						send(state.group,"reactivate", exclaimer)
					else
						send(SELF,"reactivate")
					end
				end
			end
		end]]
	>>
	
     receive damageMessage( gameObjectHandle damagingObject,
					  string damageType,
					  int damageAmount,
					  string onhit_effect )
	<<
		-- Obeliskians feel no pain!!!
		if SELF.tags.dormant then
			send(SELF,"reactivate")
		end
		
		--[[if damagingObject then 
			local attackerTags = query(damagingObject, "getTags")[1]
			if attackerTags.citizen or
				attackerTags.foreigner or
				attackerTags.bandit or
				attackerTags.animal then
				
				if state.group and damagingObject then
					send(state.group,
						"helpMemberInCombat",
						SELF, damagingObject)
				end
			end
		end]]
	>>
	
	receive emoteAffliction()
	<<

	>>
     
     receive resetEmoteTimer()
	<<
		state.AI.ints["emoteTimer"] = 0
	>>
	
	receive Vocalize( string vocalization)
	<<
		if vocalization == "idle" then
			send("rendInteractiveObjectClassHandler",
				"odinRendererPlaySFXOnInteractive",
				state.renderHandle,
				"Obeliskian Hum")
		end
	>>

	receive resetInteractions()
	<<
		printl("ai_agent", state.AI.name .. " received resetInteractions")
		
		send("rendInteractiveObjectClassHandler",
				"odinRendererClearInteractions",
				state.renderHandle)

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
		--[[
		if messagereceived == "Study Horror (obeliskian)" and
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
								"Study Horror",
								"",
								"")[1]
			end
			
			send("rendOdinCharacterClassHandler",
				"odinRendererCharacterExpression",
				state.renderHandle,
				"dialogue",
				"obeliskian",
				true)

			 send("gameBlackboard",
				"gameObjectNewJobToAssignment",
				assignment,
				SELF,
				"Study Horror (obeliskian)",
				"corpse",
				true )
			
			setCancelInteraction = true
			state.assignment = assignment
		
		-- TODO: "destroy evidence" command.
		--elseif messagereceived == "Destroy Evidence (horror)" and
			
		elseif messagereceived == "Cancel Obeliskian orders" and
			SELF.tags.dead then
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
                         "Cancel Obeliskian orders",
                         "Cancel Obeliskian orders",
                         "Cancel Obeliskian orders",
                         "Cancel Obeliskian orders",
						"",
						"",
						"Dirt",
						false)
		end]]
	>>
	
	receive AssignmentCancelledMessage( gameSimAssignmentHandle assignment )
	<<
		printl("ai_agent", state.AI.name .. " received AssignmentCancelledMessage")

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