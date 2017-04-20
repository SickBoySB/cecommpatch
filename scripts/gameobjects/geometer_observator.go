gameobject "geometer_observator" inherit "ai_agent"
<<
	local
	<<
          function geometer_observator_doOneSecondUpdate()
			state.AI.ints.secondsCount = state.AI.ints.secondsCount + 1
			if state.AI.ints.secondsCount > 100 then
				SELF.tags.exit_map = true
			end
          end
	>>

	state
	<<
		gameObjectHandle group
		gameAIAttributes AI
		gameGridPosition homeLocation
		gameGridPosition wanderDestination
		table entityData
		table traits
		string entityName
		string animSet
		int renderHandle
		int growthTimer
		bool asleep
	>>
	
	receive Create( stringstringMapHandle init )
	<<
		--local entityName = init["legacyString"]
		
		local entityName = "Geometer Observator"
          state.entityName = entityName

		local entityData = EntityDB[ entityName ]
		if not entityData then
			printl("ai_agent", "geometer_observator type not found")
			return
		end
		
		state.displayName = entityName
		
		printl("ai_agent", "placing " .. state.displayName)
		
          state.AI.strs["mood"] = "otherworldly"
		state.AI.strs["firstName"] = state.displayName
		state.AI.strs["lastName"] = ""
		state.AI.name = state.displayName -- state.AI.strs.firstName .. " " .. state.AI.strs.lastName
		state.AI.strs["gender"] = "none"
		state.AI.strs["citizenClass"] = entityName
		state.AI.ints.secondsCount = 0 
		state.AI.ints["num_afflictions"] = 0
          state.AI.ints["max_afflictions"] = 1
		 state.AI.ints["numAfflictions"] = 0 
		state.AI.ints["health"] = 15
		state.AI.ints["healthMax"] = 15
		state.AI.ints["healthTimer"] = -1
		state.AI.ints["health"] = state.AI.ints["healthMax"]
		state.AI.ints.spores = entityData.spores_per_fruit
		state.corpseTimer = 0
		
		state.group = false
		state.renderHandle = SELF.id
		
		SELF.tags = {}
		for k,v in pairs(entityData.tags ) do
			SELF.tags[v] = true	
		end
		
		state.model = entityData.model
		state.animSet = entityData.animationSet
		
		send("rendOdinCharacterClassHandler",
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
          
          state.AI.walkTicks = 3
          state.AI.ints["subGridWalkTicks"] = state.AI.walkTicks
          setposition(0,0)
	
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterWalkTicks",
			state.renderHandle,
			state.AI.walkTicks)

		local occupancyMap = entityData.occupancyMap
		local occupancyMapRotate45 = entityData.occupancyMapRotate45

		send("gameSpatialDictionary",
			"registerSpatialMapString",
			SELF,
			occupancyMap,
			occupancyMapRotate45,
			true )
		
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 0)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 1)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 2)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 3)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 4)
		
		send("gameSpatialDictionary", "gameObjectAddBit", SELF, 13) -- Geometers
		
		--send("gameSpatialDictionary", "gameObjectClearBitfield", SELF)
		--send("gameSpatialDictionary", "gameObjectClearHostileBit", SELF)
          
		--[[ doing these will crash the game btw.
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle,
			"name",
			state.entityName)
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterCustomTooltipMessage",
			SELF.id,
			"ui\\tooltips\\spectreTooltip.xml")]]
		
         
		
		--send(SELF,"makeHostile")
		ready()
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
		if SELF.tags["exited_map"] then
			send(SELF,"beDestroyed") 
			return
		end
	 
		if not state.AI or SELF.deleted then
			return
		end
		
		if SELF.tags.dead then
			state.corpseTimer = state.corpseTimer + 1
			if state.corpseTimer == 100 then
				send("rendCommandManager",
					"odinRendererCreateParticleSystemMessage",
					"QuagSmokePuff",
					state.AI.position.x,
					state.AI.position.y )
				
				send(SELF,"beDestroyed") 
			end
			
			return
		elseif state.AI.thinkLocked then
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
                    local results = query("gameBlackboard","gameAgentTestForInterruptsMessage", state.AI, SELF )
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
                         --[[send("rendOdinCharacterClassHandler",
						"odinRendererSetCharacterAttributeMessage",
                                   state.renderHandle,
                                   "currentJob",
                                   state.AI.curJobInstance.displayName)]]
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
               geometer_observator_doOneSecondUpdate()
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
			state.wanderDestination.x = state.homeLocation.x + rand(-8, 8)
			state.wanderDestination.y = state.homeLocation.y + rand(-8, 8)
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
	
     receive resetEmoteTimer()
	<<
		state.AI.ints["emoteTimer"] = 0
	>>
	
	respond getIdleAnimQueryRequest()
	<<
		local animName = "idle"
		return "idleAnimResponse", animName
	>>
	
	receive beDestroyed()
	<<
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
	
	receive deathBy( gameObjectHandle damagingObject, string damageType )
     <<
		send("rendCommandManager",
			"odinRendererCreateParticleSystemMessage",
			"QuagSmokePuff",
			state.AI.position.x,
			state.AI.position.y )

		local handle = query("scriptManager",
						 "scriptCreateGameObjectRequest",
						 "explosion",
						 { legacyString = "Small Eldritch Explosion" })[1]
		send(handle,
			"GameObjectPlace",
			state.AI.position.x,
			state.AI.position.y)
		
		send(SELF,"beDestroyed")

	>>
	
	receive damageMessage( gameObjectHandle attacker, string damageType, int damageAmount, string onhit_effect )
	<<
		SELF.tags.agitated = true
		
			if not state.AI or SELF.deleted then
			return
		end
		
		if SELF.tags.dead then
			return
		end
		
		local distX = rand(-3,3)
		local distY = rand(-3,3)

		-- teleport!
		local newLoc = gameGridPosition:new()
		local result = false
		local count = 0
		repeat  
			local dX = rand(-4,4)
			local dY = rand(-4,4)
			distX = distX + dX
			distY = distY + dY

			newLoc.x = state.AI.position.x + distX
			newLoc.y = state.AI.position.y + distY
			result = query("gameSpatialDictionary", "gridCanPathTo", SELF, newLoc, false)
			count = count + 1
		until (result == true) or (count >= 10)

		if (result == false) then
			return
		end
		
		send("rendCommandManager",
			"odinRendererCreateParticleSystemMessage",
			"QuagSmokePuffLarge",
			state.AI.position.x,
			state.AI.position.y )
		
		send(SELF, "GameObjectPlace",newLoc.x, newLoc.y)
		
		send("rendCommandManager",
			"odinRendererCreateParticleSystemMessage",
			"TransformPouf",
			state.AI.position.x,
			state.AI.position.y )
			
		send("rendInteractiveObjectClassHandler",
			"odinRendererPlaySFXOnInteractive",
			SELF.id,
			"Slipgate Return")
	>>
>>