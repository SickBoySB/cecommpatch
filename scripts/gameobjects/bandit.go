gameobject "bandit" inherit "ai_agent"
<<
	local
	<<
          function bandit_doOneSecondUpdate()
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
		gameGridPosition homeLocation
		gameGridPosition wanderDestination
		table entityData
		table traits
		table banditGroup
		string entityName
		string animSet
		int renderHandle
		bool asleep
		
	>>

	receive setGroup( gameObjectHandle group )
	<<
		-- ai_agent handles generic setGroup
		-- Here: Bandits inherit name of group when added to a group.
		
		if state.group then
			local lastName = query(state.group,"getNameSuffix")[1]
			s = string.sub(lastName,-1, #lastName)
			if s == "s" then
				lastName = string.sub(lastName,1,-2)
			end
			state.AI.strs["lastName"] = lastName
		else
			state.AI.strs["lastName"] = "Brigand"
		end
		
		state.AI.name = state.AI.strs["firstName"] .. " " ..state.AI.strs["lastName"]			
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle, "name", state.AI.name)
	>>
	
	receive Create(stringstringMapHandle init )
	<<
		local entityName = init["legacyString"]
          if not entityName then
               return
          end
          state.entityName = entityName
		
		state.group = false
		
		local entityData = EntityDB[ state.entityName ]

		if not entityData then
			printl("bandit", "bandit type not found")
			return
		end
          state.entityData = entityData

          state.AI.strs["mood"] = "covetous"
          
          if rand(0,100) < 50 then 
			state.AI.strs["gender"] = "male"
               state.AI.strs["firstName"] =  maleFirstNames[ rand(1, #maleFirstNames ) ]
		else
			state.AI.strs["gender"] = "female"
               state.AI.strs["firstName"] =  femaleFirstNames[ rand(1, #femaleFirstNames ) ]
		end
		
		state.AI.name = state.AI.strs["firstName"]
		
		local isArmoured = false
		
		local variants = {[1]="A", [2]="B", [3]="C", [4]="D"}
		state.AI.strs["variant"] = variants[rand(1,4)]
			  

		if entityName == "Armoured Bandit" then		
			state.AI.strs["citizenClass"] = "Bandit"
			isArmoured = true
			state.animSet = "biped"
			  
			send( "rendOdinCharacterClassHandler",
						"odinRendererCreateCitizen", 
						SELF, 
						entityData["torsoModel"], 
						entityData["headModel"],
						"", 
						"", 
						state.animSet, 0, 0 )
		else
			  state.AI.strs["citizenClass"] = "Bandit"
			  local models = getModelsForClass( state.AI.strs["citizenClass"],
											   state.AI.strs["gender"],
											   state.AI.strs["variant"] )

			state.animSet = models["animationSet"]

			if rand(1,3) == 1 then
				models.hatModel = entityData.hats[ rand(1,#entityData.hats)]
				send( "rendOdinCharacterClassHandler",
					"odinRendererCreateCitizen", 
					SELF, 
					models["torsoModel"], 
					models["headModel"],
					"", 
					models["hatModel"], 
					models["animationSet"], 0, 0 )
			else
			
				send( "rendOdinCharacterClassHandler",
					"odinRendererCreateCitizen", 
					SELF, 
					models["torsoModel"], 
					models["headModel"],
					models["hairModel"], 
					models["hatModel"], 
					models["animationSet"], 0, 0 )
			end
		
			state.headModel = models["headModel"]
			state.hairModel = models["hairModel"]
		end          

		send("rendOdinCharacterClassHandler", "odinRendererFaceCharacter", 
				state.renderHandle, 
				state.AI.position.orientationX,
				state.AI.position.orientationY )

		local humanstats = EntityDB["HumanStats"]
		local worldstats = EntityDB["WorldStats"]

          SELF.tags = {}
		for k,v in pairs(EntityDB.Bandit.tags) do
			SELF.tags[v] = true	
		end
		
		--[[send("gameSpatialDictionary", "gameObjectAddBit", SELF, 14) -- Bandits
		
		-- Kill 'em all
		-- empire + players 1-4
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 0)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 1)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 2)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 3)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 4)]]
          
		SELF.tags["combat_target_for_enemy"] = nil
		if isArmoured then
			SELF.tags["armoured"] = true
		end
		SELF.tags.hostile_vs_fishpeople = true
		
		state.AI.ints["grenades"] = state.AI.ints["grenadesMax"] -- and grenades, why not.

          -- START ai_damage required stats
          state.AI.ints["healthMax"] = entityData.health
		  
          state.AI.ints["healthTimer"] = 3 -- in seconds, per 1 point 
          state.AI.ints["fire_timer"] = 10
          state.AI.ints["health"] = state.AI.ints["healthMax"]
		state.AI.ints["numAfflictions"] = 0
		
          -- END ai_damage required stats
          state.AI.ints["sightRadius"] = 12
          state.AI.ints["wall_attacks"] = 0
          
          state.AI.ints["emoteTimer"] = 30
           
          state.AI.ints["morale"] = 100
          state.AI.ints["corpse_timer"] = humanstats.corpseRotTimeDays * worldstats["dayNightCycleSeconds"] * 10 -- gameticks
		state.AI.ints.corpse_vermin_spawn_time_start = div(state.AI.ints.corpse_timer,2)
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
    
          if state.AI.strs["gender"] == "female" then
			state.vocalID = "High00"
		else
			if( rand(0,100) < 50 ) then 
				state.vocalID = "Low00"
			else
				state.vocalID = "Mid00"
			end
		end

          state.AI.ints["emoteTimer"] = 0
          state.timer = 0
          
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle,
			"name",
			state.AI.name)
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterCustomTooltipMessage",
			SELF.id,
			"ui\\tooltips\\banditTooltip.xml")

          ready()
	>>

	receive gameFogOfWarExplored(int x, int y )
	<<
		state.asleep = false
	>>

	receive GameObjectPlace( int x, int y ) 
	<<
		if not state.AI.bools.first_placement then
			
			-- equip random melee & ranged weapon.
			local entity_data = EntityDB["BanditStats"]
			
			local climate = query("gameSession", "getSessionString", "biome")[1]
			
			local ranged_weapons = entity_data.ranged_weapons
			if climate == "cold" then
				ranged_weapons = entity_data.ranged_weapons_cold
			elseif climate == "desert" then
				ranged_weapons = entity_data.ranged_weapons_desert
			elseif climate == "tropics" then
				ranged_weapons = entity_data.ranged_weapons_tropics
			end
			
			send(SELF,"setWeapon","melee",entity_data.melee_weapons[ rand(1,#entity_data.melee_weapons) ] )
			send(SELF,"setWeapon","ranged",entity_data.ranged_weapons[ rand(1,#entity_data.ranged_weapons) ] )
			  
			state.homeLocation = gameGridPosition:new()
			state.homeLocation.x = x
			state.homeLocation.y = y
			
			state.AI.bools.first_placement = true
			
			send("gameSession","incSessionInt","banditsOnMap",1)
		end
		
		local banditsFaction = query("gameSession","getSessiongOH","Bandits")[1]
		if query(banditsFaction, "isHostile")[1] then
			send(SELF,"makeHostile")
		elseif  query(banditsFaction, "isFriendly")[1] then
			send(SELF,"makeFriendly")
		else
			send(SELF,"makeNeutral")
		end
		send(SELF, "resetInteractions")
	>>

	receive playerYield()
	<<
          SELF.tags["hostile_agent"] = nil
		send("gameSpatialDictionary", "gameObjectRemoveBit", SELF, 14);
          SELF.tags["peaceful"] = true
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
			"ui\\tooltips\\banditHostileTooltip.xml")
		
		send("gameSpatialDictionary","gameObjectAddBit",SELF,14)
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
			"ui\\tooltips\\banditFriendlyTooltip.xml")
		
		send("gameSpatialDictionary","gameObjectRemoveBit",SELF,14)
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
			"ui\\tooltips\\banditTooltip.xml")
		
		send("gameSpatialDictionary","gameObjectRemoveBit",SELF,14)
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
               bandit_doOneSecondUpdate()
          end
		
		if SELF.tags.friendly_agent then
			local isDay = query("gameSession","getSessionBool","isDay")[1]
			if isDay then
				send("gameSpatialDictionary",
					"gridExploreFogOfWar",
					state.AI.position.x,
					state.AI.position.y,
					state.AI.ints["sightRadius"])
			else
				-- isNight
				send("gameSpatialDictionary",
					"gridExploreFogOfWar",
					state.AI.position.x,
					state.AI.position.y,
					math.ceil(state.AI.ints["sightRadius"] * 0.5) )
			end
		end
    
		if state.AI.bools["dead"] then
			if not SELF.tags["buried"] then
				send(SELF,"corpseUpdate")
			else
				disable_buried_corpses() -- ai_agent.go function
			end	
			return
		end

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
			--return "getHomeLocationResponse", state.homeLocation
			
			-- mix it up a bit.
			state.wanderDestination.x = state.homeLocation.x + rand(-4, 4)
			state.wanderDestination.y = state.homeLocation.y + rand(-4, 4)
			return "getHomeLocationResponse", state.wanderDestination
		end
	>>
	
	receive resetInteractions()
	<<
		printl("ai_agent", state.AI.name .. " received resetInteractions")
		
		send("rendInteractiveObjectClassHandler",
			"odinRendererClearInteractions",
			state.renderHandle)
		
		if SELF.tags.dead and
			not SELF.tags.buried and
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
						 "Dump Corpses", -- "Dump Corpses",
                         "Dump Corpse (player order)", --"Dump Corpse (player order)",
						"graveyard",
						"",
						"Dirt",
						false,true)
		else
			if not SELF.tags.hostile_agent and not SELF.tags.dead then
				
				send("rendInteractiveObjectClassHandler",
						"odinRendererAddInteractions",
						state.renderHandle,
						"Shoot Bandits",
						"Shoot Bandits",
						"", --"Shoot Bandits",
						"", --"Shoot Bandits",
						"",
						"",
						"click",
						true,true)
				
			end
			
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
			state.assignment and
			SELF.tags.dead then
			send("gameBlackboard",
				"gameObjectRemoveTargetingJobs",
				SELF,
				nil)
			
			state.assignment = nil
			send(SELF,"resetInteractions")
			state.attempt_auto_corpse_handling = false
			
		elseif messagereceived == "Shoot Bandits" then
			-- it's go time!
			local banditFaction = query("gameSession","getSessiongOH","Bandits")[1]
			
			send("gameSession","setSessionInt","BanditsRelations", -45)
			send(banditFaction, "makeHostile" )
			
			send("rendCommandManager",
				"odinRendererTickerMessage",
				"You've ordered your troops to attack the Bandits, starting with " .. state.AI.name .. "!",
				"bandit_war",
				"ui\\orderIcons.xml")
			
			-- Can't undo this one!
			setCancelInteraction = false
			state.assignment = assignment
		end
		
		if setCancelInteraction then
			send("rendInteractiveObjectClassHandler",
                    "odinRendererAddInteractions",
				state.renderHandle,
                         "Cancel orders for corpse of " .. state.AI.name,
                         "Cancel corpse orders",
                         "", --"Cancel corpse orders",
                         "", --"Cancel corpse orders",
						"graveyard",
						"",
						"Dirt",
						false,true)
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
	
	receive despawn() override
	<<
		if SELF.tags["buriedandhidden"] then
			send("gameBlackboard", "gameObjectRemoveTargetingJobs", SELF, nil)
			send("rendOdinCharacterClassHandler", "odinRendererDeleteCharacterMessage", state.renderHandle)
			send("gameSpatialDictionary", "gridRemoveObject", SELF)
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
		SELF.tags.meat_source = true
		
		-- you're not a bandit anymore:
		send("gameSpatialDictionary","gameObjectRemoveBit",SELF,14) 

		local animName = bipedDeathAnimSmart(damageType) -- func in ai_agent.go
		
		if animName then
			send("rendOdinCharacterClassHandler",
				"odinRendererSetCharacterAnimationMessage",
				state.renderHandle,
				animName,
				false)
		end
		
		send("gameSession","incSessionInt","banditDeaths",1)
		send("gameSession","incSessionInt","banditsOnMap",-1)
		
		if state.AI.curJobInstance then
			FSM.abort(state,"Died.")
		end
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterCustomTooltipMessage",
			SELF.id,
			"ui\\tooltips\\banditDeadTooltip.xml")
		
		state.attempt_auto_corpse_handling = true
		send(SELF,"resetInteractions")
	>>

	receive HarvestMessage( gameObjectHandle harvester, gameSimJobInstanceHandle ji )
	<<
		SELF.tags.meat_source = nil
		send("rendCommandManager",
			"odinRendererCreateParticleSystemMessage",
			"BloodSplashCentered",
			state.AI.position.x,
			state.AI.position.y)
		
		local numSteaks = 2 -- so gross.

		for s=1, numSteaks do
				
			local results = false
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
			
			local handle = results[1]
			if( handle == nil ) then 
				printl("Creation failed")
				return "abort"
			else 
				local range = 1
				local positionResult = query("gameSpatialDictionary",
								   "nearbyEmptyGridSquare",
								   state.AI.position,
								   range)
				
				while not positionResult[1] do
					range = range + 1
					positionResult = query("gameSpatialDictionary", "nearbyEmptyGridSquare", state.AI.position, range)
				end
				if positionResult[1].onGrid then
					send( handle, "GameObjectPlace", positionResult[1].x, positionResult[1].y  )
				else
					send( handle, "GameObjectPlace", state.AI.position.x, state.AI.position.y  )
				end
				
				local civ = query("gameSpatialDictionary", "gridGetCivilization", state.AI.position )[1]
				if civ == 0 then
					send(handle,"ClaimItem")
				else
					send(handle, "ForbidItem")
				end	
			end
		end

		send( "rendOdinCharacterClassHandler",
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
		if not SELF.tags["corpse_interact"] then
			send(SELF, "resetInteractions")
			SELF.tags["corpse_interact"] = true
		end
		
		if state.AI.bools["rotted"] then
			-- now use the corpse_timer to send out detect corpse pings
			--[[
			state.AI.ints["corpse_timer"] = state.AI.ints["corpse_timer"] + 1
			
			if state.AI.ints["corpse_timer"] % 100 == 0 then -- timer in gameticks
				
				local results = query("gameSpatialDictionary",
								  "allObjectsInRadiusWithTagRequest",
								  state.AI.position,
								  10,
								  "citizen",
								  true)
				
				if results and results[1] then
					send(results[1], "hearExclamation", "detectBanditCorpse", SELF, nil)
				end
				
				
				state.AI.ints["corpse_timer"] = 1
				
			end
			]]
		else
			--broadcast that there's a rotting corpse over here.
			state.AI.ints["corpse_timer"] = state.AI.ints["corpse_timer"] - 1
			
			if state.AI.ints["corpse_timer"] % 100 == 0 then
                    -- timer in gameticks, trigger once per 10s
				
				local results = query("gameSpatialDictionary",
							 "allObjectsInRadiusWithTagRequest",
							 state.AI.position,
							 10,
							 "bandit",
							 true)
				
				if results and results[1] then
					send(results[1],"hearExclamation","detectBanditCorpse",SELF,nil)
				end
				
				if state.AI.ints["corpse_timer"] < state.AI.ints["corpse_vermin_spawn_time_start"] and
					state.numVerminSpawned < 8 then
					
					if rand(1,8) == 8 then
						local handle = query( "scriptManager",
								"scriptCreateGameObjectRequest",
								"vermin",
								{legacyString = "Tiny Beetle" } )[1]
						
						
						local positionResult = query("gameSpatialDictionary",
											    "nearbyEmptyGridSquare",
											    state.AI.position,
											    3)[1]
						send(handle,
							"GameObjectPlace",
							positionResult.x,
							positionResult.y  )

						state.numVerminSpawned = state.numVerminSpawned +1
					end
				end
				
				local civ = query("gameSpatialDictionary",
							   "gridGetCivilization",
							   state.AI.position )[1]
				
				if civ < 10 then
					if not query("gameSession",
							"getSessionBool",
							"BanditsBuryCorpsesPolicySet")[1] then
						
						local eventQ = query("gameSimEventManager",
										 "startEvent",
										 "bandits_corpses",
										 {},
										 {})[1]
						
						send(eventQ,"registerSubject",SELF)
						
					else
						if query("gameSession","getSessionBool","BanditsBuryCorpses")[1] and
							not SELF.tags.buried and
							not state.assignment and
							state.attempt_auto_corpse_handling then
							
							send(SELF,
								"HandleInteractiveMessage",
								"Bury Corpse (player order)",
								nil)
							
						elseif query("gameSession","getSessionBool","banditDumpCorpses")[1] and
							not SELF.tags.buried and
							not state.assignment and
							state.attempt_auto_corpse_handling then
							
							send(SELF,
								"HandleInteractiveMessage",
								"Dump Corpse (player order)",
								nil)
						end
					end
				end
			end

			if state.AI.ints["corpse_timer"] <= 0 then
				-- here's your skeleton model swap
				state.AI.bools["rotted"] = true
				state.AI.bools["onFire"] = false -- because we're done with that.
				SELF.tags["burning"] = false
				SELF.tags["meat_source"] = nil
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
			state.AI.ints["morale"] = state.AI.ints["morale"] - 5
			if state.AI.ints["morale"] <= 50 and not SELF.tags["fleeing"] then
				send(SELF,"moraleBroken")
			end
			
		elseif name == "detectBanditCorpse" then

			state.AI.ints["morale"] = state.AI.ints["morale"] - 3
			if state.AI.ints["morale"] <= 50 and not SELF.tags["fleeing"] then
				send(SELF,"moraleBroken")
			end	
		end
	>>
	
     receive damageMessage( gameObjectHandle damagingObject, string damageType, int damageAmount, string onhit_effect )
	<<
		if damagingObject then 
			local attackerTags = query(damagingObject, "getTags")[1]
			
			if attackerTags.animal or
				attackerTags.foreigner or
				attackerTags.obeliskian then
				
				send(damagingObject,"addTag","bandit_combat_target")
				
				if state.group and damagingObject then
					send(state.group,"helpMemberInCombat", SELF, damagingObject)
				end
			elseif attackerTags.temp_hostiles_dont_target then
				send(damagingObject,"removeTag","temp_hostiles_dont_target")
			end
			
			if not attackerTags.animal then
				-- decrease morale because getting hurt is scary
				state.AI.ints["morale"] = state.AI.ints["morale"] - damageAmount * 3
				if state.AI.ints["morale"] <= 50 and not SELF.tags["fleeing"] then
				    send(SELF,"moraleBroken")
				end
			
			end
		end
	>>
	
	receive emoteAffliction()
	<<
	     -- decrease morale after affliction because getting hurt is scary
          state.AI.ints["morale"] = state.AI.ints["morale"] - 5
          if state.AI.ints["morale"] <= 50 and not SELF.tags["fleeing"] then
              send(SELF,"moraleBroken")
          end
	>>
     
	receive moraleBroken()
	<<
		printl("ai_agent", state.AI.name .. " the Bandit received moraleBroken!")
		send(SELF,"AICancelJob", "Morale broken.")
		SELF.tags.coward = true
		if state.group then
			if SELF.tags.bandit_leader then
				if state.group then
					send(state.group,"forceRetreat") -- let's get out of here!
				end
			else
				SELF.tags.fleeing = true
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
		local spawnTable = {
			legacyString = "Bandit",
			socialClass = "lower",
			class = state.AI.strs["citizenClass"],
			gender = state.AI.strs["gender"],
			variant = state.AI.strs["variant"],
			firstName = state.AI.strs["firstName"],
			lastName = state.AI.strs["lastName"],
			numAfflictions = tostring( state.AI.ints["numAfflictions"] ), -- necessary?
		}
		
		if state.headModel then
			spawnTable.headModel = state.headModel
		end
		
		if state.hairModel then
			spawnTable.hairModel = state.hairModel
		end

		local handle = query( "scriptManager",
						"scriptCreateGameObjectRequest",
						"citizen", 
						spawnTable )[1]
		
		send(handle,
			"GameObjectPlace",
			state.AI.position.x,
			state.AI.position.y )
		
		-- send job to return to civilization?
		
		-- And now the cleanup.
		send("gameSession","incSessionInt","banditsOnMap",-1)
		if state.group then
			send(state.group, "removeMember", SELF, "joined_colony")
		end
		
		FSM.abort( state, "Joined society.")
		
		 -- This is a hack so that dead agents can't raise alarms
          SELF.tags["alarm_waypoint_active"] = true
		
		-- and remove all factional bits.
		send("gameSpatialDictionary", "gameObjectClearBitfield", SELF)
		send("gameSpatialDictionary", "gameObjectClearHostileBit", SELF)

		-- but we also need to set the waypoint flag
		send("gameSpatialDictionary", "gameObjectAddBit", SELF, 17)
		
		send(SELF,"ForceDropEverything")
		send("rendOdinCharacterClassHandler",
			"odinRendererDeleteCharacterMessage",
			state.renderHandle)
		
		send("gameSpatialDictionary",
			"gridRemoveObject",
			SELF)
		
		destroyfromjob(SELF,ji)
	>>
	
	receive DropItemMessage( int x, int y)
	<<
		send("rendInteractiveObjectClassHandler",
			"odinRendererAddInteractions",
			state.renderHandle,
			"Bury Corpse",
			"Bury Corpse (player order)",
			"",
			"",
			"body",
			"",
			"Dirt Thump", false,true)
	>>
>>