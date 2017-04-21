gameobject "geometer" inherit "ai_agent"
<<
	local
	<<
          function geometer_doOneSecondUpdate()
               --state.AI.ints["emoteTimer"] = state.AI.ints["emoteTimer"] + 1
			state.AI.ints.secondsCount = state.AI.ints.secondsCount + 1
			if state.AI.ints.secondsCount > 200 then
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
		
		local entityName = "Geometer"
          state.entityName = entityName

		local entityData = EntityDB[ entityName ]
		if not entityData then
			printl("ai_agent", "geometer type not found")
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
		state.AI.ints["health"] = 15
		state.AI.ints["healthMax"] = 15
		state.AI.ints["healthTimer"] = -1
		state.AI.ints["health"] = state.AI.ints["healthMax"]
		state.corpseTimer = 0
		
		state.group = false
		state.renderHandle = SELF.id
		
		SELF.tags = {}
		for k,v in pairs(entityData.tags ) do
			SELF.tags[v] = true	
		end
		
		local models = {}
		models.torsoModel = entityData.model
		models.headModel = entityData.headModel
		models.animationSet = entityData.animationSet

		state.models = models
		state.animSet = models["animationSet"]
		
          send("rendOdinCharacterClassHandler",
			"odinRendererCreateCitizen", 
			SELF, 
			models["torsoModel"], 
			models["headModel"],
			"", 
			"", 
			models["animationSet"],
			0,
			0 )

		send("rendOdinCharacterClassHandler",
			"odinRendererFaceCharacter", 
				state.renderHandle, 
				state.AI.position.orientationX,
				state.AI.position.orientationY )
		
          state.AI.walkTicks = entityData.walkTicks
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
		
		
		
		-- don't flag as a geometer until hostile
		if query("gameSession","getSessionBool","geometer_hostile")[1] then
			send("gameSpatialDictionary", "gameObjectAddBit", SELF, 13) -- Geometers
			
			send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 0)
			send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 1)
			send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 2)
			send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 3)
			send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 4)
			
			SELF.tags.hostile_agent = true
		end
		
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
		send("gameSpatialDictionary", "gameObjectAddBit", SELF, 13) -- Geometers
		
		SELF.tags.hostile_agent = true
		
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 0)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 1)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 2)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 3)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 4)
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
		send("gameSpatialDictionary", "gameObjectClearBitfield", SELF)
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
		send("gameSpatialDictionary", "gameObjectClearBitfield", SELF)
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
			if state.corpseTimer == 300 then
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
               local do_check = false
               if state.AI.ints.updateTimer % 3 == 0 then
                    do_check = true
               end
               
               if do_check then 
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
		
		if not state.AI or SELF.deleted then
			return
		end
		
		if state.AI.ints.updateTimer % 10 == 0 then
               geometer_doOneSecondUpdate()
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
		
		if attacker then
			local tags = query(attacker,"getTags")[1]
			if tags then
				if tags.citizen then
					if not query("gameSession","getSessionBool","geometer_hostile")[1] then
						send("gameSpatialDictionary", "gameObjectAddBit", SELF, 13) -- Geometers
					end
				end
			end
		end
	>>
	
	receive hearExclamation( string name, gameObjectHandle exclaimer, gameObjectHandle subject )
	<<
		if name == "geometer_go_hostile" then
			send(SELF,"makeHostile")
		end
	>>
>>