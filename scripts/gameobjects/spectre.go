gameobject "spectre" inherit "ai_agent"
<<
	local
	<<
          function spectre_doOneSecondUpdate()
               state.AI.ints["emoteTimer"] = state.AI.ints["emoteTimer"] + 1
			
			-- cheesy HACK til sunrise works
			if not SELF.tags.persistent then --They stay alive during the day if Persistent.
				if not query("gameSession", "getSessionBool", "isNight")[1] then
					send(SELF,"despawn")
				elseif state.owner then
					if query(state.owner,"getTags")[1].last_rites_performed then
						send(SELF,"despawn")
					elseif SELF.tags["spectre_burial"] then
						if query(state.owner,"getTags")[1].buried then
							send(SELF,"despawn")
						end
					end
				elseif SELF.tags["spectre_vengeance"] and
					SELF.tags.spectre_has_haunt_target and
					state.hauntingTarget then
					
					if query(state.hauntingTarget,"getTags").corpse then
						send(state.owner,"addTag","murder_avenged")
						send(SELF,"despawn")
					end
				end
			end
          end
	>>

	state
	<<
		gameAIAttributes AI
		gameGridPosition homeLocation
		gameGridPosition wanderDestination
		table entityData
		table traits
		string entityName
		string animSet
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
			printl("event", "spectre type not found")
			return
		end
          state.entityData = entityData

          state.AI.strs["mood"] = "spooky"
		
		if init.name then 
			state.AI.name = "The Spectre of " .. init.name
		else
			state.AI.name = "Warren Spectre"
		end
		
		SELF.tags = {}
		for k,v in pairs(entityData.tags) do
			SELF.tags[v] = true	
		end
		
		if init.tag1 then
			SELF.tags[init.tag1] = true
		end
		if init.tag2 then
			SELF.tags[init.tag2] = true
		end
		if init.goal then
			SELF.tags["spectre_" .. init.goal] = true
		else
			SELF.tags.spectre_haunting = true
		end

		state.animSet = "biped"
		
		send("rendOdinCharacterClassHandler",
			"odinRendererCreateCitizen", 
			SELF, 
			"models\\character\\body\\bipedSpectre.upm",
			"models\\character\\heads\\headSpectre.upm",
			"", 
			"", 
			state.animSet, 0, 0 )

		send("rendOdinCharacterClassHandler",
			"odinRendererFaceCharacter", 
			state.renderHandle, 
			state.AI.position.orientationX,
			state.AI.position.orientationY )
          
          state.AI.ints["emoteTimer"] = 30
		state.timer = 0
          state.AI.walkTicks = 3
          state.AI.ints["subGridWalkTicks"] = state.AI.walkTicks
          setposition(0,0)
	
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterWalkTicks",
			state.renderHandle,
			state.AI.walkTicks)

		local occupancyMap = 
			".-.\\".. 
			"-C-\\"..
			".-.\\"
		local occupancyMapRotate45 = 
			".-.\\".. 
			"-C-\\"..
			".-.\\"
			
		send("gameSpatialDictionary",
			"registerSpatialMapString",
			SELF,
			occupancyMap,
			occupancyMapRotate45,
			true )
		
		send("gameSpatialDictionary", "gameObjectClearBitfield", SELF)
		send("gameSpatialDictionary", "gameObjectClearHostileBit", SELF)
          
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle,
			"name",
			state.AI.name)
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterCustomTooltipMessage",
			SELF.id,
			"ui\\tooltips\\spectreTooltip.xml")
		
		send(SELF,"setWeapon","melee","default")
		
		state.AI.ints.morale = 100
		state.AI.ints.numAfflictions = 0
		state.AI.ints.health = 0
		
          ready()
	>>

	receive makeFriendly()
	<<
		SELF.tags.human_agent = true
	>>
	
	receive makeHostile()
	<<
		SELF.tags.hostile_spectre = true
	>>

	receive registerOwner( gameObjectHandle owner)
	<<
		state.owner = owner
	>>

	receive registerHauntingTarget( gameObjectHandle target)
	<<
		state.hauntingTarget = target
		SELF.tags.spectre_has_haunt_target = true
	>>
	
	receive makeHostile()
	<<
		--[[send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle,
			"name",
			state.AI.name)
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterCustomTooltipMessage",
			SELF.id,
			"ui\\tooltips\\banditHostileTooltip.xml")
		
		send("gameSpatialDictionary","gameObjectAddBit",SELF,14)--]]
	>>
	
	receive makeFriendly()
	<<
		--[[send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle,
			"name",
			state.AI.name)
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterCustomTooltipMessage",
			SELF.id,
			"ui\\tooltips\\banditFriendlyTooltip.xml")
		
		send("gameSpatialDictionary","gameObjectRemoveBit",SELF,14)--]]
	>>
	
	receive makeNeutral()
	<<
		--[[send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle,
			"name",
			state.AI.name)
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterCustomTooltipMessage",
			SELF.id,
			"ui\\tooltips\\banditTooltip.xml")
		
		send("gameSpatialDictionary","gameObjectRemoveBit",SELF,14)--]]
	>>
	
	receive SleepMessage()
	<<
		--state.asleep = true
	>>

     receive Update()
     <<
		if SELF.deleted then 
			return
		end
		
		if not state.AI then
			return
		end
		
  		if state.AI.thinkLocked then
               return
          end

		if state.AI.curJobInstance == nil then
               state.AI.canTestForInterrupts = true -- reset testing for interrupts 
			local results = query("gameBlackboard",
							  "gameAgentNeedsJobMessage",
							  state.AI,
							  SELF)

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
               -- interrupt only at 1s intervals because enemies getting stuff right isn't as important as humans
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
                         send("rendOdinCharacterClassHandler",
						"odinRendererSetCharacterAttributeMessage",
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
               spectre_doOneSecondUpdate()
          end
	>>
  
	receive setHomeLocation(int x, int y)
	<<
		local newLoc = gameGridPosition:new()
		newLoc.x = x
		newLoc.y = y
		state.homeLocation = newLoc
	>>
  
	respond getHomeLocation()
	<<
		if SELF.tags["idle"] then
			state.wanderDestination.x = state.homeLocation.x + rand(-10, 10)
			state.wanderDestination.y = state.homeLocation.y + rand(-10, 10)
			return "getHomeLocationResponse", state.wanderDestination
		else
			return "getHomeLocationResponse", state.homeLocation
		end
	>>
	
	receive resetInteractions()
	<<
		printl("ai_agent", state.AI.name .. " received resetInteractions")
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
		--[[local setCancelInteraction = false
		
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
			
		elseif messagereceived == "Cancel corpse orders" and
			state.assignment then
			send("gameBlackboard",
				"gameObjectRemoveTargetingJobs",
				SELF,
				nil)
			
			send(SELF,"resetInteractions")
			state.assignment = nil
			state.attempt_auto_corpse_handling = false
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
						false)
		end--]]
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

     respond getIdleAnimQueryRequest()
	<<
		-- return a random animation from predefined in characters.xml

		local idleAnims = {"idle", "idle_alt1", "idle_alt2", "idle_alt3",}
		local animName = idleAnims[rand(1,#idleAnims)]

		return "idleAnimQueryResponse", animName
	>>
	
     receive resetEmoteTimer()
	<<
		state.AI.ints["emoteTimer"] = 0
	>>
	
	respond amIDead()
	<<
		return "deadresponse", "no"
	>>
	
	receive Sunrise()
	<<
		send(SELF,"despawn")
	>>
>>