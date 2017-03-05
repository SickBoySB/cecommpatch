gameobject "selenian" inherit "ai_agent"
<<
	local
	<<
          function selenian_doOneSecondUpdate()
               --state.AI.ints["emoteTimer"] = state.AI.ints["emoteTimer"] + 1
			
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
			
			if not SELF.tags.ready_to_transform then 
				state.growthTimer = state.growthTimer + 1
				if state.growthTimer >= state.entityData.stages[ state.selenianGrowthStage ].transformTime then
					-- do transformation
					local jobname = state.entityData.stages[ state.selenianGrowthStage ].transformJobName
					
					SELF.tags.ready_to_transform = true
					
					state.growthTimer = 0
				end
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
		int selenianGrowthStage
		int growthTimer
		int corpseTimer
		bool asleep
	>>
	
	receive Create( stringstringMapHandle init )
	<<
		--local entityName = init["legacyString"]
		
		local entityName = "Selenian Polyp"
		local growthStage = 0 -- rand(0,3) --4) -- pull from init.
		
		if init.legacyString and
			(init.legacyString == "0" or
			init.legacyString == "1" or
			init.legacyString == "2" or
			init.legacyString == "3" ) then
			
			growthStage = tonumber(init.legacyString)
		end
		
		state.selenianGrowthStage = growthStage
		
		state.growthTimer = rand(1,10)
          state.entityName = entityName

		local entityData = EntityDB[ entityName ]
		if not entityData then
			printl("ai_agent", "selenian type not found")
			return
		end
		
		state.entityData = entityData
		state.displayName = entityData.stages[ state.selenianGrowthStage ].display_name
		
		printl("ai_agent", "placing " .. state.displayName)
		
          state.AI.strs["mood"] = "otherworldly"
		state.AI.strs["firstName"] = state.displayName
		state.AI.strs["lastName"] = ""
		state.AI.name = state.displayName -- state.AI.strs.firstName .. " " .. state.AI.strs.lastName
		state.AI.strs["gender"] = "none"
		state.AI.strs["citizenClass"] = entityName
		
		state.AI.ints["num_afflictions"] = 0
          state.AI.ints["max_afflictions"] = 1
		state.AI.ints["health"] = 15
		state.AI.ints["healthMax"] = 15
		state.AI.ints["healthTimer"] = -1
		state.AI.ints["health"] = state.AI.ints["healthMax"]
		state.AI.ints.spores = entityData.spores_per_fruit
		state.corpseTimer = 0
		
		state.group = false
		state.renderHandle = SELF.id
		
		SELF.tags = {}
		for k,v in pairs(entityData.stages[ state.selenianGrowthStage ].tags ) do
			SELF.tags[v] = true	
		end
		
		state.model = entityData.stages[ state.selenianGrowthStage ].model
		state.animSet = entityData.stages[ state.selenianGrowthStage ].animationSet
		
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

		local occupancyMap = ".-.\\".. 
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
		
		
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 0)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 1)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 2)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 3)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 4)
		
		send("gameSpatialDictionary", "gameObjectAddBit", SELF, 12) -- Selenians
		--send("gameSpatialDictionary", "gameObjectClearBitfield", SELF)
		--send("gameSpatialDictionary", "gameObjectClearHostileBit", SELF)
          
		--[[send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle,
			"name",
			state.entityName)
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterCustomTooltipMessage",
			SELF.id,
			"ui\\tooltips\\spectreTooltip.xml")]]
		
          ready()
		
		send(SELF,"makeHostile")
	>>
	
	receive GameObjectPlace(int x, int y)
	<<
		send("gameSession","incSessionInt","seleniansOnMap",1)
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
		if SELF.deleted or not state.AI then
			return
		end
		
		if SELF.tags.dead and state.corpseTimer then
			
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
               selenian_doOneSecondUpdate()
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
			
			local x_max = query("gameSession","getSessionInt","x_max")[1]
			local y_max = query("gameSession","getSessionInt","y_max")[1]
			
			if state.wanderDestination.x > x_max - 1 then state.wanderDestination.x = x_max - 1 end
			if state.wanderDestination.x < 1 then state.wanderDestination.x = 1 end
			
			if state.wanderDestination.y > y_max - 1 then state.wanderDestination.y = y_max - 1 end
			if state.wanderDestination.y < 1 then state.wanderDestination.y = 1 end
		
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
	
	receive despawn() override
	<<
		if state.AI.possessedObjects["curCarriedTool"] then
			send(state.AI.possessedObjects["curCarriedTool"],
				"DestroySelf",
				state.AI.curJobInstance )
		end
		
		if state.AI.possessedObjects["curPickedUpItem"] then
			send(state.AI.possessedObjects["curPickedUpItem"],
				"DestroySelf",
				state.AI.curJobInstance )
		end
		
		send("rendOdinCharacterClassHandler",
			"odinRendererDeleteCharacterMessage",
			state.renderHandle)
		
		destroyfromjob(SELF, nil )
	>>
	
     receive resetEmoteTimer()
	<<
		state.AI.ints["emoteTimer"] = 0
	>>
	
	respond getIdleAnimQueryRequest()
	<<
		local animName = "idle"
		if state.selenianGrowthStage == 1 or
			state.selenianGrowthStage== 2 then
			
			animName = "idle"
		elseif state.selenianGrowthStage == 0 then
			animName = "writhe"
		elseif state.selenianGrowthStage == 3 then
			local anims = {
				"writhe",
				"hit",
				"flip"
			}
			animName = anims[ rand(1,#anims) ]
		end

		return "idleAnimResponse", animName
	>>
	
	receive swapToNextGrowthStage()
	<<
		SELF.tags.ready_to_transform = nil
		
		-- kill old tags
		for k,v in pairs(  state.entityData.stages[ state.selenianGrowthStage ].tags ) do
			SELF.tags[v] = nil
		end
		
		state.selenianGrowthStage = state.selenianGrowthStage +1
		
		if state.selenianGrowthStage > #state.entityData.stages then
			state.selenianGrowthStage = 0
		end
		
		-- add new tags
		for k,v in pairs(  state.entityData.stages[ state.selenianGrowthStage ].tags ) do
			SELF.tags[v] = true
		end
		
		state.model = state.entityData.stages[ state.selenianGrowthStage ].model
		state.animSet = state.entityData.stages[ state.selenianGrowthStage ].animationSet
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterGeometry", 
			state.renderHandle,
			state.model, 
			"",
			"",
			state.animSet,
			"")
		
		--[[ seemingly no effect.
		send("rendOdinCharacterClassHandler",
			"odinRendererFaceCharacter", 
			state.renderHandle, 
			state.AI.position.orientationX,
			state.AI.position.orientationY)]]
		
		-- Poof!
		--[[send("rendCommandManager",
                         "odinRendererCreateParticleSystemMessage",
                         "QuagSmokePuff",
                         state.AI.position.x,
                         state.AI.position.y)]]
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
		
		if SELF.tags.drop_cool_stuff then
			local results = query( 	"scriptManager",
								"scriptCreateGameObjectRequest",
								"objectcluster",
								{ 	legacyString="Selenian Goodies",
									tagToAdd="sparkle", } )[1]
			send(results,"ClaimItem")
			send(results, "GameObjectPlace", state.AI.position.x, state.AI.position.y )
		end
			
		if state.selenianGrowthStage > 0 then
			-- non-spore selenians produce a spore upon death.
			local handle = query("scriptManager",
							 "scriptCreateGameObjectRequest",
							 "selenian",
							 { legacyString = "0" })[1]
			send(handle,
				"GameObjectPlace",
				state.AI.position.x,
				state.AI.position.y)
			
			local animName = "death"
			send("rendOdinCharacterClassHandler",
				"odinRendererSetCharacterAnimationMessage",
				state.renderHandle,
				animName,
				false)
			
			send("rendCommandManager",
					"odinRendererCreateParticleSystemMessage",
					"DustPuffLarge",
					state.AI.position.x,
					state.AI.position.y)
			
			send(SELF,"resetInteractions")
			
			SELF.tags.horror = nil
			SELF.tags.hostile_horror = nil
			SELF.tags.horror_corpse = true
		else
			local handle = query("scriptManager",
							 "scriptCreateGameObjectRequest",
							 "item",
							 { legacyString = "dormant_spore" })[1]
			send(handle,
				"GameObjectPlace",
				state.AI.position.x,
				state.AI.position.y)
			
			send("gameSession","incSessionInt","seleniansOnMap",-1)
			
			local civ = query("gameSpatialDictionary", "gridGetCivilization", state.AI.position )[1]
			if civ == 0 then
				send(handle,"ClaimItem")
			else
				send(handle,"ForbidItem")
			end	
			
			send(SELF,"beDestroyed")
		end
	>>
	
	receive damageMessage( gameObjectHandle attacker, string damageType, int damageAmount, string onhit_effect )
	<<
		if SELF.tags.dead then
			return
		end
		
		SELF.tags.agitated = true
		
		-- seriously, don't melee these things.
		
		if damageType ~= "eldritch" and
			damageType ~= "voltaic" and
			damageType ~= "explosion" and
			damageType ~= "mind_blast" and
			damageType ~= "selenian" and
			damageType ~= "fire" then
			
			local handle = query("scriptManager",
								 "scriptCreateGameObjectRequest",
								 "explosion",
								 { legacyString = "Selenian Field Effect" })[1]
				send(handle,
					"GameObjectPlace",
					state.AI.position.x,
					state.AI.position.y)
		end
	>>
>>