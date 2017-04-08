gameobject "ai_agent" inherit "renderableobject" inherit "ai_damage"
<<
	local
	<<
		function setposition( x, y )
			newPos = gameGridPosition:new()
			newPos:set( x, y )
			state.AI.position = newPos
		end
		
		function ai_agent_doOneSecondUpdate()
			state.AI.ints.emoteTimer = state.AI.ints.emoteTimer + 1
		end
	>>

	state
	<<
		gameAIAttributes AI
		int missionJobs
		int updateTimer
		int numVerminSpawned
		gameObjectHandle group
	>>

	receive Create( stringstringMapHandle init )
	<<
		initAttributes( state.AI )
		state.AI.parentGOH = SELF
		state.renderHandle = SELF.id
		
		state.missionJobs = 0
		SELF.tags.ai_agent = true
		
		state.AI.ints["emoteTimer"] = rand(1,9)
		state.AI.ints.updateTimer = rand(1,9)
		
		state.AI.strs.ranged_weapon = nil -- assume no ranged weapon by default
		SELF.tags.has_ranged_attack = nil -- maybe need this?
		
		state.AI.strs.melee_weapon = "default" -- pulls from entityDB for melee attack, or defaults to 3 blunt; if none, no melee possible
		SELF.tags.has_melee_attack = true -- maybe need this?
		
		state.AI.ints.ranged_ammo_capacity = 1
		state.AI.ints.ranged_ammo_amount = 0
		
		--state.AI.bools.has_barracks = false
		
		state.numVerminSpawned = 0
		
		state.group = nil
	>>

	receive GameObjectPlace( int x, int y ) 
	<<
		
		local x_max = 255--query("gameSession","getSessionInt","x_max")[1]
		local y_max = 255 --query("gameSession","getSessionInt","y_max")[1]
			
		if x > x_max-1 then x = x_max-1 end
		if y > y_max-1 then y = y_max-1 end
		
		if x < 1 then x = 1 end
		if y < 1 then y = 1 end
		
		local posResult = query( "gameSpatialDictionary", "gridGetObjectCenter", SELF )
		state.AI.position = posResult[ 1 ]
		
		setposition(x,y)
	
		send("gameSpatialDictionary",
			"gridMoveObjectTo", 
			SELF, 
			state.AI.position )

		send("rendOdinCharacterClassHandler", 
		     "odinRendererTeleportCharacterMessage", 
			state.renderHandle, 
			x,
			y)
		
		if SELF.tags.animal then
			state.AI.locs["herdCentre"] = state.AI.position
		end
		
		
		local iswater = query( "gameSpatialDictionary",
							"gridHasSpatialTag",
							state.AI.position,
							"water" )[1]
		
		if iswater then SELF.tags.placed_in_water = true end
	>>

	respond gridGetPosition()
	<<
		local posResult = query( "gameSpatialDictionary",
							"gridGetObjectCenter", 
							SELF )
		
		if SELF.tags.picked_up and state.carrier then
			posResult = query("gameSpatialDictionary",
							"gridGetObjectCenter", 
							state.carrier )
		end
		return "reportedPosition", posResult[1]
	>>
	
	receive setGridPosition(int x, int y)
	<<
		setposition(x,y)
	>>
	
	respond ROHQueryRequest()
	<<
		return "ROHQueryReply", state.renderHandle 
	>>
	
	respond getAIBlock()
	<<
		return "AIBlockMessage", state.AI
	>>
	
     respond getName() 
     << 
          return "nameResponse", state.AI.name
     >>
	
	respond getDisplayName()
     <<
          return "nameResponse", state.name
     >>

	respond getDisplayName() 
     << 
          return "nameResponse", state.AI.name
     >>

	receive SleepMessage()
	<<
		sleep()
		state.asleep = true
	>>
	
	receive WakeMessage()
	<<
		wake()
		if state.asleep then
			state.asleep = nil
		end
	>>

	receive makeHostile()
	<<
		if SELF.tags.dead then
			return
		end
		
		SELF.tags.hostile_agent = true
		SELF.tags.combat_target_for_enemy = nil
		SELF.tags.friendly_agent = nil
		SELF.tags.neutral_agent = nil
		
		-- NOTE: you have to have your faction flag set for hostility to work!!!
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 0)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 1)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 2)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 3)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 4)
	>>
	
	receive makeFriendly()
	<<
		if SELF.tags.dead then return end

		SELF.tags.hostile_agent = nil
		SELF.tags.combat_target_for_enemy = true
		SELF.tags.friendly_agent = true
		SELF.tags.neutral_agent = nil
		
		-- possible to be hostile to some things, however.
		send("gameSpatialDictionary", "gameObjectClearHostileBit", SELF)
	>>
	
	receive makeNeutral()
	<<
		if SELF.tags.dead then return end
	
		SELF.tags.hostile_agent = nil
		SELF.tags.combat_target_for_enemy = nil
		SELF.tags.friendly_agent = nil
		SELF.tags.neutral_agent = true
		
		send("gameSpatialDictionary", "gameObjectClearHostileBit", SELF)
	>>
	
	respond gridReportPosition()
	<<
		local posResult = query("gameSpatialDictionary",
						    "gridGetObjectCenter", 
							SELF )
		
		return "reportedPosition", posResult[ 1 ]
	>>
	
	respond getAIAttributes()
	<<
		return "AIAttributes", state.AI
	>>

	receive MoveAllowed( gameGridPosition pos )
	<<
		state.AI.bools["moveAllowed"] = true
		state.AI.position = pos
	>>

	receive MoveDenied()
	<<
		state.AI.bools["moveAllowed"] = false
	>>
	
	receive DropItemMessage( int x, int y)
	<<
		SELF.tags["picked_up"] = nil
		-- you got dropped!
		setposition(x,y)
	>>
	
	receive ForceDropEverything()
	<<
		if state.AI.possessedObjects then

			if state.AI.possessedObjects["curPickedUpItem"] and
				not state.AI.possessedObjects["curPickedUpItem"].deleted then
				
				--printl("ai_agent", state.AI.name .. " attempting to drop a thing!")
				--printl("ai_agent", " state.AI.possessedObjects[curPickedUpItem].deleted  = " .. tostring(state.AI.possessedObjects["curPickedUpItem"].deleted) )
				
				local copy = state.AI.possessedObjects["curPickedUpItem"]
			
				if not copy or copy.deleted or copy == nil then
					-- this was a problem.
					return
				end 
				
				local itemTags = query(copy, "getTags") --[1]
				if itemTags and itemTags[1] then
					itemTags = itemTags[1]
				else
					-- what is happening!
					return
				end
				local delete_me = false
				
				if itemTags.fishpeople_weapon == true or
					itemTags.despawn_on_drop == true or
					itemTags.tool == true then
					
					delete_me = true
				end
				
				if delete_me == true then
					send(copy,"DestroySelf",state.AI.curJobInstance )
				else
					-- problems with these?
					local name = query(copy,"GroundModelQueryRequest" )
					if name and name[1] then
						name = name[1]
					else
						return
					end
					
					local id = query(copy,"ROHQueryRequest")
					if id and id[1] then
						id = id[1]
					else
						return
					end
					
					send( "rendOdinCharacterClassHandler", 
							"odinRendererCharacterForceDropItemMessage",
							state.renderHandle,
							id,
							name)
					
					send(copy,
						"DropItemMessage",
						state.AI.position.x,
						state.AI.position.y)
					
					send(copy,"GameObjectPlace",
						state.AI.position.x,
						state.AI.position.y)
					
					local civ = query("gameSpatialDictionary",
								   "gridGetCivilization",
								   state.AI.position)[1]
					
					if civ < 10 then
						send(copy,"ClaimItem")
					else
						if copy and not copy.deleted then
							send(copy,"ForbidItem")
						end
					end
				end
			end

			if state.AI.possessedObjects["curPickedUpCharacter"] and
				not state.AI.possessedObjects["curPickedUpCharacter"].deleted then
				
				local copy = state.AI.possessedObjects["curPickedUpCharacter"]
				
				send("rendOdinCharacterClassHandler",
					"odinRendererCharacterDetachCharacter",
					state.renderHandle,
					value.id,
					"Bones_Group")
				
				send(copy,"DropItemMessage",state.AI.position.x, state.AI.position.y)
				send(copy,"GameObjectPlace",state.AI.position.x, state.AI.position.y)
				
				send("rendOdinCharacterClassHandler",
					"odinRendererSetCharacterAnimationMessage",
					value.id,
					"corpse_dropped",
					false)
			else
				--printl("ai_agent", "attempted to ForceDropTools but had no non-deleted contents in any state.AI.possessedObjects")
			end
		else
			--printl("ai_agent", "WARNING attempted to ForceDropTools but had no state.AI.possessedObjects")
		end
	>>

	receive RegisterItemInHandForTrade ( gameObjectHandle depot )
	<<
		if state.AI.possessedObjects["curPickedUpItem"] ~= nil then
			send(state.AI.possessedObjects["curPickedUpItem"], "RegisterItemForTrade", depot);
		end
	>>

	receive ForceCarryTradeItem ( gameObjectHandle depot, string entityName, int amount )
	<<
          if entityName == "" then
               return
          end
		if not EntityDB[ entityName ] then
               printl("ai_agent", "WARNING " .. state.AI.name .. " asked to carry invalid trade good: " .. tostring(entityName) )
			return
          end

		local resultObjects = {}
		
		for k = 1, amount do
			local handle = query( "scriptManager",
                        "scriptCreateGameObjectRequest",
                        "item",
                        { legacyString = entityName } )
			
			if handle and handle[1] then
			     resultObjects[#resultObjects + 1] = handle[1]			
				--handle[1].tags["merchant_trade_good"] = true
				
				send(handle[1],"addTag","merchant_trade_good")
				
				if state.group ~= nil then
					send(state.group, "RegisterTradeGood", handle[1])
				end
			end
		end
		
		local resultObject = nil
		
		if #resultObjects > 1 then
			-- create container, put everything into container

			local entityDBName = entityName
			local initTable = {legacyString = "crate",
						    container_parent = entityDBName }
				
			local container = query("scriptManager",
							"scriptCreateGameObjectRequest",
							"container",
							initTable )
				
			local handle = container[1]

			send(handle, "ContainerSetStackMode", entityDBName)
			for i=1,#resultObjects do
				send(handle, "ContainerAddItem", resultObjects[i])
			end
			if state.group ~= nil then
				send(state.group, "RegisterTradeGood", handle)
			end
			resultObject = handle				
		else
			resultObject = resultObjects[1]
		end
		
		if resultObject then 
			local resultROH = query(resultObject, "ROHQueryRequest" )
			local name = query( resultObject, "HandModelQueryRequest" )
			
			state.AI.possessedObjects["curPickedUpItem"] = resultObject
			
			-- Start the animation
			send( "rendOdinCharacterClassHandler", 
				  "odinRendererCharacterPickupItemMessage",
				   state.renderHandle,
				   resultROH[1],
				   "R_ItemAttach",
				   name[1],
				   "")
		end
	>>

     receive ForceCarryItem( string entityName )
     <<
          if entityName == "" then
               return
          end
          
          if state.AI.possessedObjects["curPickedUpItem"] ~= nil then
			return
		end
		
		initTable = { legacyString = entityName, }
		
		local commodityListingHideBool = "false" -- intentional string
          if SELF.tags.fishperson or
			SELF.tags.bandit or
			SELF.tags.foreigner or
			SELF.tags.temporary then
               
			initTable.hiddenFromCommodityList = "true"
          end
		
		if SELF.tags.bandit then
			-- 47B: We shall presume that all goods Bandits get spawned on them are stolen.
			initTable.tagToAdd = "stolen_goods"
		end

          local handle = query( "scriptManager",
                              "scriptCreateGameObjectRequest",
                              "item",
						initTable )[1]

          if not handle then 
               printl("ai_agent", "Force Equip Tool failed, invalid item")
               return
          end
		
		if SELF.tags.temporary then
			send(handle,"addTag","temporary")
		end
		
		local resultROH = query(handle, "ROHQueryRequest" )
		local name = query( handle, "HandModelQueryRequest" )
		
		state.AI.possessedObjects["curPickedUpItem"] = handle
          
		-- Start the animation
		send( "rendOdinCharacterClassHandler", 
			  "odinRendererCharacterPickupItemMessage",
			   state.renderHandle,
			   resultROH[1],
			   "R_ItemAttach",
			   name[1],
			   "")
     >>
     
     receive alarmWaypointReset ()
     <<
          if not SELF.tags["dead"] then
               SELF.tags["alarm_waypoint_active"] = nil
			send("gameSpatialDictionary", "gameObjectRemoveBit", SELF, 17)
          end
     >>
     
     receive deathBy( gameObjectHandle damagingObject, string damageType )
     <<
		-- clear attached assignment for burial/whatever.
		state.assignment = nil
		
          -- This is a hack so that dead agents can't raise alarms
          SELF.tags["alarm_waypoint_active"] = true
		SELF.tags.marked_for_beating = nil
		-- but we also need to set the waypoint flag
		send("gameSpatialDictionary", "gameObjectAddBit", SELF, 17)
		
		-- and remove all factional bits.
		send("gameSpatialDictionary", "gameObjectClearBitfield", SELF)
		send("gameSpatialDictionary", "gameObjectClearHostileBit", SELF)
	    
		state.AI.bools.dead = true
		SELF.tags.corpse = true
		SELF.tags.dead = true
		SELF.tags.hostile_agent = nil
		SELF.tags.friendly_agent = nil
		SELF.tags.combat_target_for_enemy = nil -- will become deprecated post 43? -- no. 
		SELF.tags["conversable"] = nil
		SELF.tags.occult_inspector_destroy_target = nil
		
		send( "gameSpatialDictionary",
			"registerSpatialMapString",
			SELF, "c", "c", true )
	    
		if state.group ~= nil and
			state.group ~= false and
			state.group and
			not state.group.deleted then
			
			if SELF.tags["cultist"] then
				send("gameSession", "incSessionInt", "cultPower", -1)
			end
			
			local reason = "died"
			if damagingObject ~= nil then
				
				local killerTags = query(damagingObject,"getTags")[1]
				
				if killerTags.fishperson then reason = "fishperson"
				elseif killerTags.bandit then reason = "bandit"
				elseif killerTags.novorusian then reason = "novorusian"
				elseif killerTags.mecharepublicain then reason = "mecharepublicain"
				elseif killerTags.stahlmarkian then reason = "stahlmarkian"
				elseif killerTags.obeliskian then reason = "obeliskian"
				elseif killerTags.citizen and killerTags.military then reason = "military"
				elseif killerTags.animal then reason = "animal"
				end
			end
						
			send(state.group,"removeMember", SELF, reason)
		end
		
		send(SELF,"putOutFire")
		--send(SELF,"resetInteractions") -- do this per subclass due to order of operations re. how data is cleaned up
     >>

	receive wakeFromSleep()
	<<
		if state.AI.bools["asleep"] then
			-- sleep.fsm will detect when this bool is flipped and do the correct wakeup sequence
			state.AI.bools["asleep"] = false
		end
	>>
	
	respond getTags()
	<<
		return "getTagsResponse", SELF.tags
	>>

	receive hearExclamation( string name, gameObjectHandle exclaimer, gameObjectHandle subject )
	<<
		if SELF.tags.dead then return end
		
		if name == "explosion" then
			send(SELF, "wakeFromSleep")
		elseif name == "gunshot" then
			send(SELF, "wakeFromSleep")
		elseif name == "loudnoise" then
			send(SELF, "wakeFromSleep")
		end
	>>
	
	receive attemptEmote( string emoteIcon, int emoteTimerMin, bool any_zoom )
	<<
		if state.AI.ints["emoteTimer"] > emoteTimerMin then
			send("rendOdinCharacterClassHandler",
					"odinRendererCharacterExpression",
					state.renderHandle,
					"thought",
					emoteIcon,
					any_zoom )
			
			send(SELF, "resetEmoteTimer")
		end
	>>
	
	receive AICancelJob(string reason)
	<<
		if state.AI.curJobInstance then
			printl("ai_agent", state.AI.name .. " cancelled job: " .. state.AI.curJobInstance.name .. ", because: " .. tostring(reason))
			FSM.abort( state, reason )
			state.AI.curJobInstance = nil
		end
	>>
	
	receive Update()
	<<
		if SELF.tags.destroy_me then
			destroyfromjob(SELF, nil)
			return
		end
	
		if SELF.tags.dead then
			return
		end
		
		--DET_REC( "Script Agent " .. tostring( SELF ) .. " being updated" )

		-- Do a 'sanity' check to make sure that this character makes spatial sense
		-- The types of things that are impassable for this agent
		
		local impassability = {"landscape", "object"}

		if not SELF.tags[ "amphibious" ] then 
			impassability[ #impassability + 1 ] = "water"
		end

		local allowedBlockers = {}

		if SELF.tags[ "amphibious" ] then 
			allowedBlockers[ #allowedBlockers + 1 ] = "water"
		end		

		if SELF.tags[ "human"]  then 
			allowedBlockers[ #allowedBlockers + 1 ] = "door"
		end	

		local pathresults = query( "gameSpatialDictionary", 
								"gridSanityCheck",  
								SELF,
		       	     		 	state.AI.position, 
								impassability,
		       	     		 	allowedBlockers )

		if pathresults.name == "gridPathStuck" then
			printl("ai_agent", "Failed sanity check: " .. tostring( SELF ) .. ": stuck at " .. tostring( state.AI.position ) )
			-- if we're stuck, we want to find a new place to go 
			-- this really is error recovery -- so we'll just teleport 
			local emptyarearesults = query( "gameSpatialDictionary", 
							       	   	    "nearbyEmptyAreaNear",  
							       	   		state.AI.position,
							       	   		SELF,
							       	   		impassability )

			send( "rendOdinCharacterClassHandler", 
				  "odinRendererTeleportCharacterMessage", 
				  state.renderHandle, 
				  emptyarearesults[ 1 ].x,
				  emptyarearesults[ 1 ].y )

			state.AI.position = emptyarearesults[ 1 ]							   
			printl("ai_agent", tostring( SELF ) .. ": now at " .. tostring( state.AI.position ) )
		end		

		-- end spatial sanity check

		state.AI.ints.updateTimer = state.AI.ints.updateTimer +1
		if state.AI.ints.updateTimer % 10 == 0 then
			
			state.AI.ints["emoteTimer"] = state.AI.ints["emoteTimer"] +1
			
			if SELF.tags.burning and not SELF.tags.fire_immune then
				send(SELF,"damageMessage", nil,"fire",1,nil)
			end
		end
		
		-- can be killed from fire damage above, check before doing updatetimer.
		if not SELF.tags.dead and
			state.AI.ints.updateTimer % 30 == 0 then
			-- 3 second update.
			state.AI.ints.updateTimer = 0
		end
	>>
	
	receive resetEmoteTimer()
	<<
		state.AI.ints["emoteTimer"] = 0
	>>
	
	receive addTag( string name )
	<<
		SELF.tags[name] = true
	>>
	
     respond getTags()
     <<
          return "getTagsResponse", SELF.tags
     >>
	
	receive removeTag( string name )
	<<
		SELF.tags[name] = nil
	>>
	
	receive thinkLockMessage( gameObjectHandle thinkLocker, bool x)
	<<
		state.AI.thinkLocker = thinkLocker
		state.AI.thinkLocked = x
	>>

	receive setCanBeSocial(bool value)
	<<
		state.AI.bools["canBeSocial"] = value
		SELF.tags["conversable"] = value
	>>

	respond getCanBeSocial()
	<<
		if state.AI.thinkLocked then
			return "canBeSocialResponse", false
		end

		return "canBeSocialResponse", state.AI.bools["canBeSocial"]
	>>

	respond isSitting()
	<<
		return "sittingResponse", state.AI.bools["sitting"]
	>>
	
	respond isHoldingItem()
	<<
		if state.AI.possessedObjects then
			for key, value in pairs(state.AI.possessedObjects) do
				if (key == "curPickedUpCharacter" and value) or
					(key == "curPickedUpItem" and value) then
					
					return "isHoldingItemResponse", true
				end
			end
		end
		return "isHoldingItemResponse", false
	>>
	
	respond getGroup()
	<<
		if state.group then
			return "getGroupResponse", state.group
		else
			return "getGroupResponse", nil
		end
	>>
	
	receive setGroup(gameObjectHandle group)
	<<
		printl("ai_agent", "received setGroup from group named: " .. query(group,"getName")[1] )
		state.group = group
		if group then
			send(group,"addMember",SELF)
		end
		if type(group) == "cult_group" then
			printl("ai_agent", "group is a cult")
			-- set character as "inACult" so the renderer can know
			state.AI.ints["inACult"] = 1
		end
	>>
	
	receive setAIStringAttribute ( string name, string value)
	<<
		state.AI.strs[name] = value
	>>
	
	receive setAIIntAttribute ( string name, int value)
	<<
		state.AI.ints[name] = value
	>>

	receive setAIBoolAttribute ( string name, int value)
	<<
		state.AI.bools[name] = value
	>>
	
	respond getAIStringAttribute(string name)
	<<
		return "getAIBoolAttributeMessage", state.AI.strings["name"]
	>>
	
	respond getAIIntAttribute(string name)
	<<
		return "getAIBoolAttributeMessage", state.AI.ints["name"]
	>>
	
	respond getAIBoolAttribute(string name)
	<<
		return "getAIBoolAttributeMessage", state.AI.bools["name"]
	>>
	
	receive removeHostileBit(int bit)
	<<
		send("gameSpatialDictionary", "gameObjectRemoveHostileBit", SELF, bit)
	>>
	
	receive addHostileBit(int bit)
	<<
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, bit)
	>>
	
	receive setCarrier( gameObjectHandle holder)
	<<
		if holder then
			state.carrier = holder
		else
			state.carrier = nil
		end
	>>
	
	receive setRelatedDirectorName( string director_name)
	<<
		state.AI.strs.director_name = director_name
	>>
	
	receive setWeapon( string weaponType, string weaponEntityName)
	<<
		printl("ai_agent", state.AI.name .. " received setWeapon: " ..
			  tostring(weaponType) ..
			  " / " .. tostring(weaponEntityName) )
		
		if weaponType == "melee" then
			if not weaponEntityName or weaponEntityName == "none" then
				state.AI.strs.melee_weapon = nil
				SELF.tags.has_melee_attack = nil
			elseif weaponEntityName == "default" then
				state.AI.strs.melee_weapon = "default"
				SELF.tags.has_melee_attack = true
			else
				local weaponInfo = EntityDB[ weaponEntityName ]
				state.AI.strs.melee_weapon = weaponEntityName
				SELF.tags.has_melee_attack = true
			end
			
		elseif weaponType == "ranged" then
			if not weaponEntityName or weaponEntityName == "none" then
				state.AI.ints.ranged_ammo_capacity = 1
				state.AI.ints.ranged_ammo_amount = 0
				state.AI.strs.ranged_weapon = nil
				SELF.tags.has_ranged_attack = nil
			elseif weaponEntityName == "default" then
				state.AI.ints.ranged_ammo_capacity = 1
				state.AI.ints.ranged_ammo_amount = 0
				state.AI.strs.ranged_weapon = "default"
				SELF.tags.has_ranged_attack = true
			else
				local weaponInfo = EntityDB[ weaponEntityName ]
				state.AI.strs.ranged_weapon = weaponEntityName
				state.AI.ints.ranged_ammo_capacity = weaponInfo.ammo_capacity
				state.AI.ints.ranged_ammo_amount = 0 -- weaponInfo.ammo_capacity
				SELF.tags.has_ranged_attack = true
			end
		end
	>>
	
	respond getPronoun()
	<<
		if state.AI.strs["gender"] == "male" then
			return "pronounResponse", "he","him","his","He","Him","His"
		elseif state.AI.strs["gender"] == "female" then
			return "pronounResponse", "she","her","her","She","Her","Her"
		else
			return "pronounResponse", "they","them","their","They","Their","Them"
		end
	>>
	
	receive placeInGrave( gameObjectHandle burier)
	<<
		local newPos = gameGridPosition:new()
		newPos.x = state.AI.position.x
		newPos.y = state.AI.position.y
		state.gravePosition = newPos
		
		SELF.tags.buried = true
		
		printl("ai_agent", state.AI.name ..
			  " was buried in a grave at " .. tostring(state.AI.position.x) .. " / " .. tostring(state.AI.position.y))
	>>
	
	respond getGravePosition()
	<<
		if state.gravePosition then
			return "getGravePositionResponse", state.gravePosition
		end
		return "getGravePositionResponse", nil
	>>

	receive despawn()
	<<
		FSM.abort( state, "Despawning.")
		if SELF.tags.spectre then
			send("rendCommandManager",
					"odinRendererCreateParticleSystemMessage",
					"QuagSmokePuffLarge",
					state.AI.position.x,
					state.AI.position.y )
		
			-- The jig is up!
			printl("ai_agent", state.AI.name .. " being removed from game.")
			
			
			state.AI.bools["canBeSocial"] = false
			send("rendOdinCharacterClassHandler",
				"odinRendererHideCharacterMessage",
				state.renderHandle,
				true)
			
			send("rendOdinCharacterClassHandler",
				"odinRendererDeleteCharacterMessage",
				state.renderHandle)
			
			send("gameSpatialDictionary",
				"gridRemoveObject",
				SELF)
			
			send("gameBlackboard",
				"gameObjectRemoveTargetingJobs",
				SELF,
				nil)
			
			--send(SELF,"AICancelJob", "despawning")
			send(SELF,"ForceDropEverything")
			SELF.tags.destroy_me = true
			return
		end
		
		send("rendOdinCharacterClassHandler", "odinRendererDeleteCharacterMessage", state.renderHandle)
          send("gameSpatialDictionary", "gridRemoveObject", SELF)
          send("gameBlackboard", "gameObjectRemoveTargetingJobs", SELF, nil)
		destroyfromjob(SELF,nil)
	>>
	
	receive registerTradeOffice( gameObjectHandle trade_office )
	<<
		state.trade_office = trade_office
	>>
>>