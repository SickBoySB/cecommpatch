gameobject "item" inherit "object_damage"
<<
	local 
	<<

		function HideCommodityListHelper( b )
			SELF:setTemporarilyHidden(b)
		end
	
		function createStockpileJob()
            -- stops people from returning fishpeople weapons to stockpiles
            -- ... until such a point where we WANT to do this.

            if not SELF.tags["fishpeople_weapon"] then
				send("gameBlackboard", "AddObjectToStockpileList", SELF, state.parent.filter, state.EntityDBName );
				state.inStockpileList = true
            end	
		end

		function MaybeSetupStockpileJob()
			if SELF.tags["in_container"] then
				return
			end

			local results = query("gameZoneManager", "SquareInStockpile", state.position.x, state.position.y)
			if results then
				if results[1] and state.parent and state.parent.filter then
					filterResult = query(results[1], "stockpileHasFilter", state.parent.filter)
					if filterResult and not filterResult[1] then
						createStockpileJob()
					else
						state.inStockpileList = false
					end
				elseif not results[1] then
					createStockpileJob()
				end
			end
		end

		function resetItemUse()
			local forbidden = state.forbidden
			local for_trade = state.flagged_for_trade
			
			printl("item", tostring(SELF.id) .. " resetting: forbidden: " .. tostring(forbidden) .. ", for_trade: " .. tostring(for_trade))
			
			send("rendInteractiveObjectClassHandler",
				"odinRendererClearInteractions",
				state.renderHandle)

			if forbidden then
				send("rendStaticPropClassHandler",
					"odinRendererTintStaticProp",
					state.renderHandle, 255, 0, 0, 64)
				
				send("rendStaticPropClassHandler",
					"odinRendererSetStaticPropOutlineFlags", state.renderHandle, true, false, false);
				
--[[
			elseif for_trade then
				send("rendStaticPropClassHandler",
					"odinRendererTintStaticProp",
					state.renderHandle, 0, 255, 0, 64)
				
				send("rendStaticPropClassHandler",
					"odinRendererSetStaticPropOutlineFlags", state.renderHandle, false, true, false);
--]]				
			else
				send("rendStaticPropClassHandler",
					"odinRendererTintStaticProp",
					state.renderHandle, 255, 255, 255, 0)
				
				send("rendStaticPropClassHandler",
					"odinRendererSetStaticPropOutlineFlags",
					state.renderHandle, false, false, false)
				
			end

			if forbidden then			
				send("rendInteractiveObjectClassHandler",
					"odinRendererAddInteractions",
					state.renderHandle,
					"Claim Item",
					"Claim Item",
					    "Claim Item", --"Claim Item",
					    "Claim Item", --"Claim Item",
					    "chop_tree_icon",
							"",
							"chop_wood",
							false,true)
			else
				send("rendInteractiveObjectClassHandler",
					"odinRendererAddInteractions",
					state.renderHandle,
						"Forbid Item",
						"Forbid Item",
						"", --"Forbid Item",
						"", --"Forbid Item",
						"chop_tree_icon",
							"",
							"chop_wood",
							false,true)
					    
				    -- the icon indicates that it is claimed
				    
--[[
				if state.parent.trade_good == true then
					--If you're not a trade good, no trade options for you.
					if for_trade then
						send("rendInteractiveObjectClassHandler",
						    "odinRendererAddInteractions",
						    state.renderHandle,
						    "Cancel Trade Designation",
						    "Cancel Trade Designation",
							    "", --"Cancel Trade Designation",
							    "", --"Cancel Trade Designation",
							    "chop_tree_icon",
								    "",
								    "chop_wood",
								    false,true)
					else
						send("rendInteractiveObjectClassHandler",
						    "odinRendererAddInteractions",
						    state.renderHandle,
						    "Designate As Trade Good",
						    "Designate As Trade Good",
							    "", --"Designate As Trade Good",
							    "", -- "Designate As Trade Good",
							    "chop_tree_icon",
								    "",
								    "chop_wood",
								    false,true)
					end
				end
--]]
			end
		end
	>>

	state
	<<
		int renderHandle
		int amount

		bool forbidden
		bool locked_for_future_use
		bool flagged_for_trade
		
		gameGridPosition position
		gameGridPosition stashedGGP
		table parent
		string EntityDBName
		string name
		string attachedPosition
		bool inStockpileList
		int usesRemaining
		bool in_container
		
		gameObjectHandle parent_container
	>>

	receive Create( stringstringMapHandle init )
	<<
		local name = init["legacyString"]
          printl("item", "creating a " .. name)
		
          state.amount = 1 -- setting up for future stuff
          if name == "Crate of Combat Supplies" then -- hax
               state.amount = 100
          end
               
          if init["amount"] then
               state.amount = tonumber( init["amount"] )
			if state.amount == 0 then
				ScriptError("item", "ITEM CREATION: attempting to make " .. name .. " with amount=0")
			end
          end
       
		state.locked_for_future_use = false        
		state.in_container = false
		state.inStockpileList = false
		state.parent_container = nil
		state.EntityDBName = name
		
		state.forbidden = false			-- probably?
		state.flagged_for_trade = false

		state.name = name
		state.parent = EntityDB[name]
          if state.parent.display_name then
               state.displayName = state.parent.display_name
          else
               state.displayName = state.name
          end
		state.usesRemaining = 1 -- start loaded
          
          if init.displayNameOverride then
               state.displayName = init.displayNameOverride
          end
          
		if state.parent == nil then
			ScriptError("item", "ITEM CREATION: Couldn't find " .. name)
		end

		SELF.tags = state.parent.tags
		SELF.tags.commodity = true
		SELF:setCommodityName(state.parent.name)

		SELF.tags["can_receive_goods"] = true

		if init.tagToAdd then
			SELF.tags[init.tagToAdd] = true
		end

		state.renderHandle = SELF.id

          
          if state.parent.ground_models then
               state.ground_model = state.parent.ground_models[ rand(1, #state.parent.ground_models) ]
          else
               state.ground_model = state.parent.ground_model
          end
          
          if state.parent.hand_models then
               state.hand_model = state.parent.hand_models[ rand(1, #state.parent.hand_models) ] 
          else
               state.hand_model = state.parent.hand_model
          end
          
  		send("rendStaticPropClassHandler",
			"odinRendererCreateStaticPropRequest", 
			SELF, 
			state.ground_model,
			50, 50)
          
          if SELF.tags["tinned_exotic_caviar"] then
               send("gameSession", "incSessionInt", "tipTopCaviarCount", 1)
          end
        
		local description = ""
		if state.parent and state.parent.description then
			description = state.parent.description
		end


		local gameplayDescription = ""
		if state.parent and state.parent.gameplay_description then
			gameplayDescription = state.parent.gameplay_description
		end

		local tradeValue = -1
		if state.parent and state.parent.trade_value then
			tradeValue = state.parent.trade_value
		end

		send("rendInteractiveObjectClassHandler",
               "odinRendererBindCommodityTooltip",
               state.renderHandle,
               "ui//tooltips//commodityTooltipDetailed.xml",
               state.name,
               description,
			   tradeValue,
			   gameplayDescription,
			   1,
			   nullHandle)
          
--[[
        if init.hiddenFromCommodityList == "true" then
            -- yes, a string.
            -- don't do IncCommodity; for fishweapons esp.
			state.hiddenFromCommodityList = true
		elseif SELF.tags.hiddenFromCommodityList then
			state.hiddenFromCommodityList = true
        else
            send("gameUtilitySingleton", "IncCommodity", state.parent.name, 1)
        end
		--]]
		
		send( "gameSpatialDictionary", "registerSpatialMapString", SELF, "c", "c", true )
		resetItemUse()

		ready()
		sleep()
		
		if name == "cult_juice" or
			name == "fish_juice" or
			name == "happy_juice" or
			name == "anger_juice" then
			
			send("gameBlackboard",
				"gameObjectNewJobByNameMessage",
				SELF,
				"Drink Me",
				"drink" ) 
		end
	>>
	
	receive CreatedStockpileInThisSquare()
	<<
		local results = query("gameZoneManager", "SquareInStockpile", state.position.x, state.position.y)
		if results then
			if results[1] and state.parent and state.parent.filter then
				filterResult = query(results[1], "stockpileHasFilter", state.parent.filter)
				if filterResult and not filterResult[1] then
					createStockpileJob()
				else
					send("gameBlackboard", "RemoveObjectFromStockpileList", SELF)
					state.inStockpileList = false
				end
			end
		end
	>>

	receive AddedToCollection ( gameObjectHandle collection )
	<<
	--[[
		if collection == getLocalPlayerHandle() then
			if not state.hiddenFromCommodityList and not state.in_commodity_list then
				send("gameUtilitySingleton", "IncCommodity", state.parent.name, 1)
				state.in_commodity_list = true
			end
		end
		--]]
	>>

	receive RemovedFromCollection ( gameObjectHandle collection )
	<<
	--[[
		if collection == getLocalPlayerHandle() then
			if not state.hiddenFromCommodityList and state.in_commodity_list then
				printl("CONSOLE: dec commodity " .. state.parent.name .. " from RemovedFromCollection")
	            send("gameUtilitySingleton", "DecCommodity", state.parent.name, 1)
				state.in_commodity_list = false
			end
		end
		--]]
	>>

	receive setHiddenFromCommodityList( bool b )
	<<
		HideCommodityListHelper(b);
	>>

	receive StockpileRemoved()
	<<
		if state.inStockpileList == false then
			createStockpileJob()
		else

		end
	>>


	receive AddedToContainer( gameObjectHandle container )
	<<
		state.in_container = true
		SELF.tags["in_container"] = true
		state.parent_container = container
		send("rendStaticPropClassHandler", "odinRendererSetStaticPropHiddenMessage", SELF.id, true)
		if state.inStockpileList == true then
			send("gameBlackboard", "RemoveObjectFromStockpileList", SELF);
			state.inStockpileList = false
		end

		send("gameContainerManager", "RemoveItem", goh);
	>>

	receive RemovedFromContainerNoVis( gameObjectHandle container )
	<<
		state.in_container = false
		SELF.tags["in_container"] = nil
		state.parent_container = nil		
	>>

	receive RemovedFromContainer( gameObjectHandle container )
	<<
		state.in_container = false
		SELF.tags["in_container"] = nil
		state.parent_container = nil
		
		send("rendStaticPropClassHandler", "odinRendererSetStaticPropHiddenMessage", SELF.id, false)
	>>

	receive ForceRemoveContainer()
	<<
		if state.in_container == true then		
			send(state.parent_container, "ContainerRemoveItem", SELF, nil)
			SELF.tags["in_container"] = nil
			state.parent_container = nil
			send("rendStaticPropClassHandler", "odinRendererSetStaticPropHiddenMessage", SELF.id, false)
		end

	>>
	respond IsInContainer()
	<<
		return "InContainerResponse", state.in_container
	>>

	respond SmartObjectCanDoBehaviour ( string name, stringVector objectTags, stringVector notTags, gameObjectHandle owner, gameSimJobInstanceHandle ji )
	<<
		if name == "Return Goods" then
			if state.in_container == true then
				return "SmartObjectCannotDoBehaviourResponse"
			end
--			if owner:hasCommonOwner(SELF) then
				local respo = query(owner, "GetItemParentRequest")
				if respo and respo[1] == state.EntityDBName then
					return "SmartObjectCanDoBehaviourResponse", SELF, state.position, true, false
				end
--			end
		end

		return "SmartObjectCannotDoBehaviourResponse"
	>>

	respond SmartObjectRequestBehaviour (string name )
	<<
		if name == "Return Goods" then
			return "SmartObjectBehaviourResponse", "Create Container"
		else
			return "SmartObjectCannotDoBehaviourResponse"
		end
	>>

     receive ClearViolently( gameSimJobInstanceHandle ji, gameObjectHandle damagingObject )
     <<
          SELF.tags["destroyed"] = true
          if SELF.tags["explosive"] then
               if SELF.tags["crate_of_ammunition"] then
                    local results = query("scriptManager",
								"scriptCreateGameObjectRequest",
								"explosion",
								{ legacyString="Ammo Explosion" } )
				
                    if results and results[1] then
                         send(results[1],
						"GameObjectPlace",
						state.position.x,
						state.position.y )
                    end 
               else
                    local results = query("scriptManager",
								"scriptCreateGameObjectRequest",
								"explosion",
								{ legacyString="Medium Explosion" } )
				
                    if results and results[1] then
                         send(results[1],
						"GameObjectPlace",
						state.position.x,
						state.position.y )
                    end
			end
		else
			--this is temporary until we have a better way to handle it.
			send("rendCommandManager",
				"odinRendererCreateParticleSystemMessage",
				"TreeChoppedDown",
				state.position.x, state.position.y)
			
			send("rendInteractiveObjectClassHandler",
				"odinRendererPlaySFXOnInteractive",
				state.renderHandle,
				"Break (Dredmor)" ) --i hate this object so much
          end

          send(SELF,"DestroySelf", ji )
     >>

	receive HideCommodity()
	<<
		HideCommodityListHelper(true)
		SELF.tags["hidden"] = true
		send("rendStaticPropClassHandler", "odinRendererDeleteStaticProp", SELF.id)
		send("gameSpatialDictionary", "gridRemoveObject", SELF)
		send("gameContainerManager", "RemoveItem", SELF);
	>>
	
	receive UnhideCommodity( gameGridPosition newpos )
	<<
		HideCommodityListHelper(false)
		SELF.tags["hidden"] = nil
		send(SELF, "GameObjectPlace", newpos.x, newpos.y)
		send("gameContainerManager", "AddItem", SELF, newpos)
		send("gameContainerManager", "SetItemParent", SELF, state.name)
	>>
	
	receive DestroyedMessage()
	<<
		--[[
        if state.in_commodity_list then
			if not state.hiddenFromCommodityList then
				printl("CONSOLE: dec commodity " .. state.parent.name .. " from DestroyedMessage")
				send("gameUtilitySingleton", "DecCommodity", state.parent.name, 1)
			end
		end
		--]]

		if state.in_container == true then
			send(state.parent_container, "ContainerRemoveItem", SELF, ji );
		else
			send("gameContainerManager", "RemoveItem", SELF);
		end

		send("gameTradeGoodManager", "RemoveTradeGood", state.parent.name, SELF);
		if SELF.tags["merchant_trade_good"] == true then
		
		end

		send("rendStaticPropClassHandler", "odinRendererDeleteStaticProp", SELF.id);
		send("gameTradeGoodManager", "RemoveOtherTradeGood", state.parent.name, SELF);
		send("gameBlackboard", "RemoveObjectFromStockpileList", SELF);
		send("gameContainerManager", "RemoveItem", SELF);
	>>

	receive DestroyFromModule( gameSimJobInstanceHandle ji )
	<<
		send(SELF,"DestroySelf", ji )
	>>

	receive ReleaseFromModule() override
	<<
		if state.in_container and state.in_container == true then
			-- note: we should never hit this

		else
			send("rendStaticPropClassHandler", "odinRendererSetStaticPropHiddenMessage", SELF.id, false)
			send(SELF, "itemDroppedMessage", state.stashedGGP)
		end
	>>

	receive StashedMessage( gameSimJobInstanceHandle ji,  gameGridPosition ggp)
	<<
		if state.in_container and state.in_container == true then
			-- note: we should never hit this
		else
			state.stashedGGP = ggo
			send("rendStaticPropClassHandler", "odinRendererSetStaticPropHiddenMessage", SELF.id, true);
		end
		if ji then
			ji.moduleStoredGoods[#ji.moduleStoredGoods+1] = SELF
		end
	>>

	receive BuildingLocked()
	<<
		if not state.locked_for_future_use then
			state.locked_for_future_use = true
			HideCommodityListHelper(true)
		end
	>>

	receive BuildingUnlocked()
	<<
		if state.locked_for_future_use then
			state.locked_for_future_use = false
			HideCommodityListHelper(false)
			send(SELF, "GameObjectPlace", state.position.x, state.position.y)
		end
	>>

	receive ForbidItem()
	<<
		-- If I am a trade good, I will no longer be a trade good.
		
		state.flagged_for_trade = false
		SELF.tags["trade_good"] = false
		if not state.forbidden and isOwnedByLocalPlayer(SELF) then
			if SELF.tags.food then
				local foodCount = query("gameSession", "getSessionInt", "foodCount")[1]
				if foodCount then
					if foodCount <= 0 then
						send("gameSession", "setSessionInt", "foodCount", 0)
					else
						foodCount = foodCount - 1
						send("gameSession", "setSessionInt", "foodCount", foodCount)
					end
				end
			end
			if SELF.tags.cooked then
				local cookedFoodCount = query("gameSession", "getSessionInt", "cookedFoodCount")[1]
				if cookedFoodCount then
					if cookedFoodCount <= 0 then
						send("gameSession", "setSessionInt", "cookedFoodCount", 0)
					else
						cookedFoodCount = cookedFoodCount - 1
						send("gameSession", "setSessionInt", "cookedFoodCount", cookedFoodCount)
					end
				end
			end
		end

		send("gameTradeGoodManager", "RemoveTradeGood", state.parent.name, SELF)
		
		state.forbidden = true
		clearOwnedByLocalPlayer(SELF)	-- also handles the commodity list.
		resetItemUse()
	>>

	receive ClaimItem()
	<<
		state.forbidden = false		

		if not isOwnedByLocalPlayer(SELF) then --to prevent double adds
			if SELF.tags.food then
				local foodCount = query("gameSession", "getSessionInt", "foodCount")[1]
				if foodCount then
					foodCount = foodCount + 1
					send("gameSession", "setSessionInt", "foodCount", foodCount)
				else
					send("gameSession", "setSessionInt", "foodCount", 1)
				end
			end
			if SELF.tags.cooked then
				local cookedFoodCount = query("gameSession", "getSessionInt", "cookedFoodCount")[1]
				if cookedFoodCount then
					cookedFoodCount = cookedFoodCount + 1
					send("gameSession", "setSessionInt", "cookedFoodCount", cookedFoodCount)
				else
					send("gameSession", "setSessionInt", "cookedFoodCount", 1)
				end
			end
		end


		setOwnedByLocalPlayer(SELF) -- note: this now handles the commodity list
		send(SELF, "SetForTradeGood") -- everything is a trade good!
		resetItemUse()
	>>

	receive SetForTradeGood()
	<<
		if not isOwnedByLocalPlayer(SELF) then
			return
		end

		if state.parent.trade_good then 
			if not state.flagged_for_trade then
				state.flagged_for_trade = true
				SELF.tags["trade_good"] = true
				send("gameTradeGoodManager", "AddTradeGood", state.parent.name, SELF)
				resetItemUse()
			end
		end
	>>

	receive SetForNotTradeGood()
	<<
		if not isOwnedByLocalPlayer(SELF) then
			return
		end

		if state.parent.trade_good then 
			if state.flagged_for_trade == true then
				
				state.flagged_for_trade = false
				SELF.tags["trade_good"] = false
				send("gameTradeGoodManager", "RemoveTradeGood", state.parent.name, SELF)
			end
			resetItemUse()
		end
	>>

	receive InteractiveMessage( string messagereceived )
	<<
		printl ("item", tostring(SELF.id) .. " / " .. tostring(state.parent.name) .. " got InteractiveMessage: " .. messagereceived )

		if messagereceived == "Forbid Item" then
			send(SELF,"ForbidItem")
			
		elseif messagereceived == "Claim Item" then
			send(SELF,"ClaimItem")
			
--[[
		elseif messagereceived == "Designate As Trade Good" then 
			send(SELF, "SetForTradeGood")
			
		elseif messagereceived == "Cancel Trade Designation" then
			send(SELF, "SetForNotTradeGood")
--]]			
		elseif messagereceived == "Destroy Item (test)" then
               
               send("rendStaticPropClassHandler", "odinRendererDeleteStaticProp", SELF.id)
               send("gameSpatialDictionary", "gridRemoveObject", SELF)
               send(SELF, "DestroyedMessage")
               destroy(SELF)
          end
	>>

	receive RegisterItemForTrade ( gameObjectHandle depot )
	<<			
		send("rendCommandManager", "IncOtherMerchantTradeCount", depot, state.EntityDBName, 1);
	>>

	receive GameObjectPlace(int x, int y)
	<<
		state.position.x = x
		state.position.y = y
		send("gameSpatialDictionary", "gridAddObjectTo", SELF, state.position);
		send("gameContainerManager", "AddItem", SELF, state.position);		-- for the container case, the container then removes the item again
		send("gameContainerManager", "SetItemParent", SELF, state.name);
		send("rendStaticPropClassHandler", "odinRendererMoveStaticProp", SELF.id, state.position.x, state.position.y);

		MaybeSetupStockpileJob()		
	>>

	receive GameObjectPlaceContainerAssignment( int x, int y, gameSimAssignmentHandle assignment )
	<<
		state.position.x = x
		state.position.y = y
		send("gameSpatialDictionary", "gridAddObjectTo", SELF, state.position);
		send("rendStaticPropClassHandler", "odinRendererMoveStaticProp", SELF.id, state.position.x, state.position.y);
		
		results = query("gameZoneManager", "SquareInStockpile", x, y)
		if results then
			if results[1] and state.parent and state.parent.filter then
				filterResult = query(results[1], "stockpileHasFilter", state.parent.filter)
				if filterResult and not filterResult[1] then
					createStockpileJob()
				else
					state.inStockpileList = false
				end
			elseif not results[1] then
				createStockpileJob()
			end
		end

		--[[
		state.position.x = x
		state.position.y = y
		send("gameSpatialDictionary", "gridAddObjectTo", SELF, state.position);
		send("rendStaticPropClassHandler", "odinRendererMoveStaticProp", SELF.id, state.position.x, state.position.y);
        send( "gameBlackboard",
                    "gameObjectNewJobToAssignment",
                    assignment,
                    SELF,
                    "Return Goods To Assignment Container",
                    "goods",
                    true )
		--]]
	>>
	
	respond getUsesRemaining()
	<<
		return "usesRemainingResponse", state.usesRemaining
	>>
	
	respond gridReportPosition()
	<<
		return "gridReportedPosition", state.position 
	>>

	respond gridGetPosition()
	<<
		return "reportedPosition", state.position
	>>
  
	respond getItemFilter()
	<<
		if state.parent and state.parent.filter then
			return "itemFilterResponse", state.parent.filter
		else
			return "itemFilterResponse", "";
		end
	>>
  
     respond getName()
     <<
          return "getNameResponse", state.name
     >>
	
	respond getDisplayName()
     <<
          return "getNameResponse", state.displayName
     >>
	
	receive addTag( string tag )
	<<
		SELF.tags[tag] = true
	>>
	
	receive removeTag( string tag )
	<<
		SELF.tags[tag] = nil
	>>
	
	respond checkTag( string tag )
	<<
		if SELF.tags[tag] then
			return "checkTagResponse", true
		else
			return "checkTagResponse", false
		end
	>>
	
     respond getTags()
     <<
          return "getTagsResponse", SELF.tags
     >>
  
	receive itemPickedUpMessage( gameSimJobInstanceHandle ji )
	<<
		printl("item", "itemPickedUpMessage :  ROH " .. state.renderHandle);
		if state.inStockpileList == true  then
			printl("item","Removing due to stockpile job.");
			send("gameBlackboard", "RemoveObjectFromStockpileList", SELF)
			state.inStockpileList = false
		else
			printl("item","No stockpile job.");
		end
	
		if state.in_container == true then
			send(state.parent_container, "ContainerRemoveItem", SELF, ji );
		else
			send("gameContainerManager", "RemoveItem", SELF);
		end
	>>

	receive itemDroppedMessage( gameGridPosition ggp )
	<<
		if not ggp then
			printl("item", "WARNING! nil ggp given. Game will now Make Something Up.")
			local startX = query("gameSession", "getSessionInt", "startX")[1]
			local startY = query("gameSession", "getSessionInt", "startY")[1]
			
			ggp = gameGridPosition:new()
			ggp.x = startX + rand(-4,4)
			ggp.y = startY + rand(-4,4)
		end
		
		send("gameSpatialDictionary","gridAddObjectTo",SELF,ggp)
		send("gameContainerManager", "AddItem", SELF, ggp);		-- for the container case, the container then removes the item again
		send("gameContainerManager", "SetItemParent", SELF, state.name)
	>>
	
	receive stockpileFilterChanged(string name, bool value)
	<<
		if state.parent and state.parent.filter == name then
			if value then
				if state.inStockpileList then
					send("gameBlackboard", "RemoveObjectFromStockpileList", this)					
					state.inStockpileList = false;
				end
			else
				-- is now handled somewhere else
			end
		end
	>>

	receive Update()
	<<

	>>

	respond UseCommodityAnimationSet()
	<<
		return "UseCommodityAnimationSetResult", true
	>>

	respond useAnimation()
	<<
		return "useAnimationReply", state.parent.use_animation
	>>
	
		respond startAnimation()
	<<
		return "useAnimationReply", state.parent.start_animation
	>>

	respond loopAnimation()
	<<
		return "useAnimationReply", state.parent.loop_animation
	>>

	respond endAnimation()
	<<
		return "useAnimationReply", state.parent.end_animation
	>>

	receive BoardVehicle ( gameObjectHandle gOH )
	<<
		printl("item", "boarding vehicle GOH");
		local resultROH = query( gOH, "ROHQueryRequest" );
		printl("item", "boarding vehicle GOH: got a resultROH "..resultROH[1]);
		seatingArrangement = query ( gOH, "cargoRequest" );
		printl("item", "BoardVehicle: Got seat " .. seatingArrangement[1]);

		send("rendOdinCharacterClassHandler", "odinRendererCharacterPickupItemMessage", resultROH[1], state.renderHandle, seatingArrangement[1], "");
		send( gOH, "AttachedCargo", SELF );
		state.attachedPosition = seatingArrangement[1]
	>>

	receive UnboardVehicle ( gameObjectHandle gOH, int offsetX, int offsetY )
	<<
		newPosition = query(gOH, "gridGetPosition");
		newPosition[1].x = newPosition[1].x + offsetX
		newPosition[1].y = newPosition[1].y + offsetY
		localresultROH = query( gOH, "ROHQueryRequest" );
		send("rendOdinCharacterClassHandler", "odinRendererCharacterDropItemMessage", resultROH[1], state.renderHandle, "", newPosition[1].x, newPosition[1].y, "", newPosition[1].x, newPosition[1].y);
		state.position.x = newPosition[1].x
		state.position.y = newPosition[1].y
		printl("item", "Unboard vehicle: moving item to " .. newPosition[1].x .. ", " .. newPosition[1].y);
		
		send("gameSpatialDictionary", "gridAddObjectTo", SELF, state.position);
		send("rendStaticPropClassHandler", "odinRendererMoveStaticProp", SELF.id, state.position.x, state.position.y);	

		send( "gameBlackboard", "gameObjectNewJobByNameMessage", SELF, "Return Goods", "goods" )
	>>

	respond ROHQueryRequest()
	<<
		--printl("item", "ROH Query request, returning render handle " .. state.renderHandle)
		return "ROHQueryReply", state.renderHandle 
	>>

	respond GroundModelQueryRequest()
	<<
		--printl("item", "Ground Model Query Request, returning " .. state.parent.ground_model)
		return "groundModelName", state.ground_model 
	>>

	respond HandModelQueryRequest()
	<<
		--printl("item", "Hand Model Query Request returning " .. state.parent.hand_model)
		return "handModelName", state.hand_model 
	>>

	respond GetItemParentRequest()
	<<
		return "itemReportedParent", state.EntityDBName
	>>


	
	receive Transform( gameObjectHandle worker, string commodity, gameSimJobInstanceHandle ji )
	<<
		-- Put some music in the dynamics.
		incMusic(3,10);

		local results = query( "scriptManager", "scriptCreateGameObjectRequest", "item", {legacyString = commodity} );
				
			handle = results[1];
			if( handle == nil ) then 
				printl("Creation failed")
				return "abort"
			else
				send( handle, "GameObjectPlace", state.position.x, state.position.y  )
				send("rendStaticPropClassHandler", "odinRendererRotateStaticProp", handle.id , rand(0, 359), 0.25);
				-- job to return object
			end
		-- Finaeely, destroy myself.
		send (SELF, "DestroySelf", ji )

	>>
  
	receive DestroySelf( gameSimJobInstanceHandle ji )
	<<
		send(SELF,"ForbidItem")
		send(SELF, "DestroyedMessage" )
		send("rendStaticPropClassHandler", "odinRendererDeleteStaticProp", SELF.id)
		send("gameSpatialDictionary", "gridRemoveObject", SELF)
		destroyfromjob(SELF, ji)
	>>
	
	-- are these two receivers used anywhere?
	receive gameObjectAddedToWorld (gameGridPosition g)
	<<
		send("rendCommandManager", "commodityAddedToPosition", state.EntityDBName, g.x, g.y);
	>>
	
	receive gameObjectRemovedFromWorld (gameGridPosition g)
	<<
		send("rendCommandManager", "commodityRemovedFromPosition", state.EntityDBName, g.x, g.y);
	>>
  
     respond GetName()
     <<
          if state.EntityDBName then
               return "GetNameResponse", state.EntityDBName
          else
               return "GetNameResponse", "no name"
          end
     >>

	receive DropItemMessage( int x, int y )
	<<
		--printl("item","DROPPING ITEM: " .. tostring(SELF.id))
		state.position.x = x
		state.position.y = y
		
		if SELF.tags.fishpeople_weapon then
			send(SELF, "DestroySelf", nil )
		end

		MaybeSetupStockpileJob()
	>>

	receive markForDestruction()
	<<
		send("gameBlackboard",
			"gameObjectRemoveTargetingJobs",
			SELF,
			nil)
		
		send("rendInteractiveObjectClassHandler",
			"odinRendererClearInteractions",
			state.renderhandle)
		
		send("rendCommandManager",
			"odinRendererCreateParticleSystemMessage",
			"Small Beacon",
			state.position.x,
			state.position.y)
		
		local assignment = query("gameBlackboard",
							"gameObjectNewAssignmentMessage",
							SELF,
							"Destroy Item",
							"",
							"")[1]
		
		send("gameBlackboard",
			"gameObjectNewJobToAssignment",
			assignment,
			SELF,
			"Destroy Item (order)",
			"item",
			true)
		
		--state.assignment = assignment
	>>
>>