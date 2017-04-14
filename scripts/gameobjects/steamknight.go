gameobject "steamknight" inherit "ai_agent"
<<
	local
	<<
		function getSteamknightLastName()
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
	
			lastname = string.upper( lastname:sub(1,1) ) .. lastname:sub(2, #lastname)
			return lastname
		end
		
		function getSteamknightFirstName()
			local ranks = {
				"Cadet",
				"Lieutenant",
				"Captain",
				"Major",
			}
			
			return ranks[ rand(1,#ranks ) ]
		end
		
		function steamknight_doOneSecondUpdate()

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

				if state.buildingTimer > 30 then 
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
	>>

	state
	<<
		gameObjectHandle group
		gameObjectHandle attach_target
		gameAIAttributes AI
		gameGridPosition spawnLocation
		gameGridPosition wanderDestination
		table entityData
		string entityName
		string animSet
		int renderHandle
		int timer
		int buildingTimer
		bool asleep
		--bool occupied	
		--gameObjectHandle occupant			
		--int disabledRemovalTimer
	>>

	receive Create( stringstringMapHandle init )
	<<
		local entityName = "Steam Knight" --init["legacyString"]
		if not entityName then
			printl("ai_agent", "steamknight name not found: " .. tostring(entityName))
			return
		end
		state.entityName = entityName
		local ED = "Steam Knight"
		if init.skType then
			ED = init.skType
		end
		local entityData = EntityDB[ED]
		--local entityData = EntityDB[ state.entityName ]
		if not entityData then
			printl("ai_agent", "steamknight type not found")
			return
		end
		
		if init.firstName then
			state.AI.strs["firstName"] = init.firstName
		else
			state.AI.strs["firstName"] = getSteamknightFirstName()
		end
		
		if init.lastName then
			state.AI.strs["lastName"] = init.lastName
		else
			state.AI.strs["lastName"] = getSteamknightLastName()
		end
		
		if init.departureDays then
			if init.departureDays == "-1" then
				state.departureCounter = -1
			else
				state.departureCounter = tonumber(init.departureDays) * EntityDB.WorldStats.dayNightCycleTenthSeconds
			end
		else
			state.departureCounter = 10000 -- 2 * EntityDB.WorldStats.dayNightCycleTenthSeconds
		end
		
		--state.AI.strs["firstName"] = ranks[ rand(1,#ranks) ]
		--state.AI.strs["lastName"] = lastName
		
		state.AI.name = state.AI.strs["firstName"] .. " " .. state.AI.strs["lastName"]
		
		state.AI.strs["gender"] = "none"
		
		state.AI.strs["citizenClass"] = entityName
		state.model = entityData.model
		--state.animSet = entityData.animationSet
		state.animSet = entityData.animationSet
		 
		--[[
		-- CECOMMPATCH - disabling, this is causing a model to exist at 0,0 for some reason
		send("rendOdinCharacterClassHandler",
			"odinRendererCreateCharacter", 
			SELF,
			state.model,
			state.animSet,
			0,
			0 )
			]]--
		
		send("rendOdinCharacterClassHandler",
			"odinRendererFaceCharacter", 
			state.renderHandle, 
			state.AI.position.orientationX,
			state.AI.position.orientationY )
          
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
		end
		
		
		 -- START ai_damage required stats
		if entityData.health then
               state.AI.ints.healthMax = entityData.health.healthMax
               state.AI.ints.health = entityData.health.healthMax
               state.AI.ints.healthTimer = 0
               state.AI.ints.healthTimerMax = entityData.health.healthTimerSeconds
		else
			state.AI.ints.healthMax = 20
               state.AI.ints.health = 20
               state.AI.ints.healthTimer = -1
               state.AI.ints.healthTimerMax = 3
          end
		state.AI.ints["numAfflictions"] = 0
          --state.AI.ints["fire_timer"] = 10
		
		--[[ state.AI.bools["disabled"] = true
          state.renderHandle = SELF.id
          state.occupied = false]]
          
         
          
          -- START ai_damage required stats
          --[[state.AI.ints["healthMax"] = 20
          state.AI.ints["healthTimer"] = -1 -- never regenerate! Requires repairs. 
          state.AI.bools["active_melee_combat"] = false
          state.AI.ints["health"] = state.AI.ints["healthMax"]]
		
		local colour = entityData.colours[rand(1,#entityData.colours)]
		state.colour = colour 
		state.modelHead = entityData.models[ colour ].head[ rand(1,#entityData.models[ colour ].head )]
		state.modelBody = entityData.models[ colour ].body[ rand(1,#entityData.models[ colour ].body )]

		send("rendOdinCharacterClassHandler",
			"odinRendererCreateCitizen", 
			SELF,
			--entityData.model_unloaded,
			--entityData.models[ rand(1,#entityData.models) ],
			--entityData.model_head,
			state.modelBody,
			state.modelHead,
			"",
			"",
			entityData.animationSet,
			--entityData.animationSet_unloaded,
			0,
			0 )
		
		state.AI.walkTicks = entityData.walkTicks
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterWalkTicks", 
			SELF.id,
			state.AI.walkTicks )
		
		state.AI.ints["subGridWalkTicks"] = state.AI.walkTicks

		-- tell the engine that I'm a steam knight

		send("rendOdinCharacterClassHandler","odinRendererSetCharacterSightRadius", SELF.id, 20 )

		--state.AI.strs["loadout_tool"] = "grenade_launcher"
		if entityData.ranged_weapons then
			send(SELF,"setWeapon","ranged", entityData.ranged_weapons[ rand(1,#entityData.ranged_weapons) ])
		end
		send(SELF,"setWeapon","melee", entityData.melee_weapons[ rand(1,#entityData.melee_weapons) ])
		
		state.AI.strs["occupancyMap"] = 
			"....@....\\"..
			"..-----..\\".. 
			"@-**C**-@\\"..
			"..-----..\\"..
			"....@....\\"

		state.AI.strs["occupancyMapRotate45"] =  
			"....@....\\"..
			"..-----..\\".. 
			"@-**C**-@\\"..
			"..-----..\\"..
			"....@....\\"

		send( "gameSpatialDictionary",
			"registerSpatialMapString", 
			SELF, 
			state.AI.strs["occupancyMap"], 
			state.AI.strs["occupancyMapRotate45"], 
			false )

		send("gameSpatialDictionary", "gameObjectClearBitfield", SELF)
		send("gameSpatialDictionary", "gameObjectClearHostileBit", SELF)
		
		-- I am a subject of the empire and player 1, so bits 0 and 4
		send("gameSpatialDictionary", "gameObjectAddBit", SELF, 0)
		send("gameSpatialDictionary", "gameObjectAddBit", SELF, 4)

		-- I am hostile to all other players.
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 1)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 2)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 3)

		-- I am hostile to things with these bits set:
		if query("gameSession", "getSessionInt", "RepubliqueRelations")[1] <
			query("gameSession", "getSessionInt", "RepubliqueNeutralHostile")[1] then
			
			send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 5) -- Republique
		end
		
		if query("gameSession", "getSessionInt", "StahlmarkRelations")[1] <
			query("gameSession", "getSessionInt", "StahlmarkNeutralHostile")[1] then
			
			send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 6) -- Stahlmark
		end
		
		if query("gameSession", "getSessionInt", "NovorusRelations")[1] <
			query("gameSession", "getSessionInt", "NovorusNeutralHostile")[1] then
			
			send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 7) -- Novorus
		end
		
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 8) -- Carnivores
		-- 9 = herbivores
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 10) -- Fishpeople
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 11) -- Obeliskians
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 12) -- Selenians
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 13) -- Geometers
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 14) -- Bandits
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 15) -- Frontier Justice Targets
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 16) -- Cultist Murder Targets
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle,
			"name",
			state.AI.name)
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterCustomTooltipMessage",
			SELF.id,
			"ui\\tooltips\\steamknightTooltip.xml")
		
		state.AI.ints["morale"] = 100 -- or whatever?
		
		state.AI.bools["keep_temp_weapon"] = true
		
		state.attach_target = nil
		
		-- start hidden, reveal via drop
		send("rendOdinCharacterClassHandler",
			"odinRendererHideCharacterMessage",
			state.renderHandle,
			true)
		
		send("gameBlackboard",
			"gameCitizenJobToMailboxMessage",
			SELF,
			nil,
			"Drop In (steamknight)",
			"")
		
		-- add self to collection.
		local collection = query("gameObjectManager", "gameObjectCollectionRequest", "steamKnights")[1]
		collection[ #collection +1 ] = SELF
		send("gameSession","incSessionInt","steamKnightsActive",1)
		ready()
	>>

	receive GameObjectPlace( int x, int y ) 
	<<
		-- look before we drop. Don't land in structures.
		
		
		
		--[[local isInvalidDrop = true
		local i = 0
		while isInvalidDrop do
			if i > 0 then
				state.AI.position.x = x + rand(i * -1,i)
				state.AI.position.y = y + rand(i * -1,i)
			end
			local iswater = query("gameSpatialDictionary","gridHasSpatialTag",state.AI.position,"water" )[1]
			
			
			local isstructure = query("gameSpatialDictionary",
							"gridHasSpatialTag",
							state.AI.position,
							"occupiedByStructure" )[1]
			
			if not iswater and not isstructure then
				isInvalidDrop = false
			end
			
			i = i + 1
		end]]
		
		local newLoc = gameGridPosition:new()
		newLoc.x = x
		newLoc.y = y
		
		local positionResult = query("gameSpatialDictionary", "nearbyEmptyGridSquare", newLoc, 10)[1]
		
		setposition(x,y)

		send("rendOdinCharacterClassHandler", 
				"odinRendererTeleportCharacterMessage", 
				state.renderHandle, 
				x,
				y)
		
		state.spawnLocation = gameGridPosition:new()
		state.spawnLocation.x = x
		state.spawnLocation.y = y
		
		
		send("gameSpatialDictionary",
			"gridExploreFogOfWar",
			state.spawnLocation.x,
			state.spawnLocation.y,
			10)
		
		--send(SELF,"resetInteractions")
	>>
	
	respond seatingRequest()
	<<
		-- DEPRECATED
		printl("Got a seating request!");
		return "SeatingReply", "Chest"
	>>

	receive BoardVehicle2 ( gameObjectHandle passenger, int rOH, string seat )
	<<
		-- DEPRECATED
		send("rendOdinCharacterClassHandler",
			"odinRendererCharacterPickupCharacterMessage", 
			state.renderHandle, 
			rOH, 
			seat,
			"board_steamknight", 
			"models\\vehicles\\steamKnight.upm", 
			"steamknight", 
			"close",
			"stand",
			true)
    
		state.occupant = passenger
		state.AI.bools["disabled"] = false
		state.disabledRemovalTimer = 50
	>>

	receive AttachedCargo( gameObjectHandle cargo )
	<<
		-- DEPRECATED
	>>

	receive CharacterBoarding ( gameObjectHandle passenger )
	<<
		-- DEPRECATED
		send(passenger, "BoardVehicle", SELF);
	>>

	receive Update()
	<<
		tooltip_refresh_from_save()
		
		if state.AI.thinkLocked then
               return
          end
		
		if state.AI.bools["dead"] or state.AI.bools["disabled"] then
			if not SELF.tags["sk_exploded"] then
				if state.AI.ints["corpse_timer"] then
					state.AI.ints["corpse_timer"] = state.AI.ints["corpse_timer"] - 1
					
					if state.AI.ints["corpse_timer"] <= 0 then
						-- makes sense for SKs to explode, right?
						-- yes, this is the second explosion trigger (first in deathBy).. it's intentional
						local results = query("scriptManager",
									"scriptCreateGameObjectRequest",
									"explosion",
									{ legacyString="Ammo Explosion" } )
					
						if results and results[1] then
							 send(results[1],
							"GameObjectPlace",
							state.AI.position.x,
							state.AI.position.y )
						end 
						
						SELF.tags["sk_exploded"] = true
						
						--send(SELF,"despawn") -- disabling for now... can't figure out how to get rid of the weapon
					end
				else
					state.AI.ints["corpse_timer"] = 20
				end
			end
			
			return
		end
		
		-- do a deathcheck here?
		if state.departureCounter then
			if state.departureCounter ~= -1 then
				if state.departureCounter <= 0 then
					-- queue job to leave.
					SELF.tags.exit_map = true
					
				elseif state.departureCounter == 600 then
					-- give warning.
					send("rendCommandManager",
						"odinRendererStubMessage",
						"ui\\orderIcons.xml",
						"icon_steamknight",
						"Steam Knight Leaves Soon", -- header text
						state.AI.name .. " the Steam Knight will be leaving very soon.", -- text description
						"Right-click to dismiss.", -- action string
						"steamKnightDeparted", -- alert type (for stacking)
						"ui//eventart//battle.png", -- imagename for bg
						"high", -- importance: low / high / critical
						SELF.id, -- object ID
						45 * 1000, -- duration in ms
						0, -- "snooze" time if triggered multiple times in rapid succession
						nil) -- gameobjecthandle of director, null if none
				end
				state.departureCounter = state.departureCounter - 1
			end
		end
		
		send("gameSpatialDictionary",
			"gridExploreFogOfWar",
			state.AI.position.x,
			state.AI.position.y,18)
		
          if state.AI.ints.updateTimer % 10 == 0 then
               steamknight_doOneSecondUpdate()
          end
    
		send("rendOdinCharacterClassHandler",
			"odinRendererCharacterSetIntAttributeMessage",
			SELF.id,
			"morale",
			state.AI.ints["morale"])
		
		send("rendOdinCharacterClassHandler",
			"odinRendererCharacterSetIntAttributeMessage",
			SELF.id,
			"health",
			state.AI.ints["health"])
		
		send("rendOdinCharacterClassHandler",
			"odinRendererCharacterSetIntAttributeMessage",
			SELF.id,
			"healthMax",
			state.AI.ints["healthMax"])
		
		send("rendOdinCharacterClassHandler",
			"odinRendererCharacterSetIntAttributeMessage",
			SELF.id,
			"numAfflictions",
			state.AI.ints["numAfflictions"])
		
		if state.AI.curJobInstance == nil then
               state.AI.canTestForInterrupts = true -- reset testing for interrupts 
			local results = query("gameBlackboard",
							  "gameAgentNeedsJobMessage",
							  state.AI,
							  SELF )

			if results.name == "gameAgentAssignedJobMessage" then 
				-- CECOMMPATCH - hacky temp fix for alarm performance spikes
				if results[1].displayName == "Responding To Alarm" then
					return
				end
				
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
               
               --if oneSecond then
                    local results = query("gameBlackboard",
								  "gameAgentTestForInterruptsMessage",
								  state.AI,
								  SELF )
                    if results.name == "gameAgentAssignedJobMessage" then

						-- CECOMMPATCH - hacky temp fix for alarm performance spikes
						if results[1].displayName == "Responding To Alarm" then
							return
						end

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
			--[[else
				local keepStepping = true
				while keepStepping do
					keepStepping = FSM.step(state) 
				end		
			end]]
          end
	>>

	receive InteractiveMessage( string messagereceived )
	<<

	>>

	receive MoveAllowed(gameGridPosition pos)
	<<
		state.AI.bools["moveAllowed"] = true
		state.AI.position = pos

		send("gameSpatialDictionary",
			"gridExploreFogOfWar",
			state.AI.position.x,
			state.AI.position.y,16)		-- FIXME: sightradius
	>>

	receive MoveDenied()
	<<
		state.AI.bools["moveAllowed"] = false
	>>

	respond isVehicleEmpty()
	<<
		-- DEPRECATED
		return "VehicleStatus", state.occupied
	>>
  
	receive deathBy( gameObjectHandle damagingObject, string damageType )
	<<		
		-- makes sense for SKs to explode, right?
		local results = query("scriptManager",
					"scriptCreateGameObjectRequest",
					"explosion",
					{ legacyString="Medium Explosion" } )
	
		if results and results[1] then
			 send(results[1],
			"GameObjectPlace",
			state.AI.position.x,
			state.AI.position.y )
		end 
		
		-- remove self from collection
		local collection = query("gameObjectManager", "gameObjectCollectionRequest", "steamKnights")[1]
		for k, v in pairs(collection) do
			if v == SELF then
				collection[k] = nil
				break
			end
		end
		
		send("gameSession","incSessionInt","steamKnightsActive",-1)
		
		state.AI.bools["dead"] = true
		SELF.tags["dead"] = true
		
		--s = state.AI.strs["firstName"] .. " " .. state.AI.strs["lastName"] .. " has died due to " .. damageType .. "!"
		--send("rendCommandManager", "odinRendererTickerMessage", "A Steamknight has been killed!", "skull", "ui\\thoughtIcons.xml")
		--send(SELF, "Vocalize", "Dying")
		
		
		-- animation time
		local randomDeath = rand(1, 2)

		if randomDeath == 1 then
			send("rendOdinCharacterClassHandler",
				"odinRendererSetCharacterAnimationMessage",
				state.renderHandle,
				"die_front", false)
		else --randomDeath == 2 then
			send("rendOdinCharacterClassHandler",
				"odinRendererSetCharacterAnimationMessage",
				state.renderHandle,
				"die_back", false);
		end
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterCustomTooltipMessage",
			SELF.id,
			"ui\\tooltips\\steamknightDeadTooltip.xml")
				
		--SELF.tags["human"] = false
		--SELF.tags["corpse"] = true
		incMusic(4,30); -- ???
		incMusic(5,5);
		--send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "nameprefix", "The Late ");
		printl("ai_agent", "Steamknight successfully died!")
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
			"ui\\tooltips\\steamknightTooltip.xml")
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
			"ui\\tooltips\\steamknightTooltip.xml")
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
			"ui\\tooltips\\steamknightTooltip.xml")
	>>
	
	receive attachToCharacter( gameObjectHandle target )
	<<
		--[[printl("ai_agent", "attaching " .. state.AI.name .. " to target: " .. query(target,"getName")[1] )
		state.attach_target = target
		SELF.tags.attached = true
		
		local name = query(state.attach_target,"getName")[1]
		send("rendCommandManager",
			"odinRendererStubMessage",
			"ui\\orderIcons.xml", -- iconskin
			"nco_command_icon", -- icon
			"Steam Knight Assigned", -- header text
			state.AI.name .. " the Steam Knight has been attached to " .. name .. "'s squad.", -- text description
			"Left-click to zoom. Right-click to dismiss", --"Left-click to zoom. Right-click to dismiss.", -- action string
			"", -- alert type (for stacking)
			"ui//eventart//battle.png", -- imagename for bg
			"low", -- importance: low / high / critical
			SELF.id, --query(handle2, "ROHQueryRequest")[1], -- object ID
			45 * 1000, -- duration in ms
			0, -- "snooze" time if triggered multiple times in rapid succession
			nil)]]
	>>
	
	receive despawn() override
	<<
		-- remove me from world, take carried objects with me.
		-- if carrying a body (for some reason), leave it behind.
		
		send("gameBlackboard", "gameObjectRemoveTargetingJobs", SELF, nil)
		send("rendOdinCharacterClassHandler", "odinRendererDeleteCharacterMessage", state.renderHandle)
		send("gameSpatialDictionary", "gridRemoveObject", SELF)
		destroyfromjob(SELF,ji)
	>>
>>