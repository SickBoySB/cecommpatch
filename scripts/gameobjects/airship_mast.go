gameobject "airship_mast" inherit "object_damage" inherit "spatialobject"
<<
	local 
	<<
		function airship_mast_reset_controls()
			
			send("rendInteractiveObjectClassHandler",
				"odinRendererClearInteractions",
				state.renderHandle)
			
			if not SELF.tags["under_construction"] then
				if SELF.tags["active_airship_mast"] then
					
					send("rendInteractiveObjectClassHandler",
						"odinRendererAddInteractions",
						state.renderHandle,
						"De-activate Airship Signaling",
						"De-activate Airship Signaling",
						"", --"De-activate Airship Signaling",
						"", --"De-activate Airship Signaling",
						"",
						"",
						"click01",
						false,true)
					
				else
					-- if not active_airship_mast
					
					send("rendInteractiveObjectClassHandler",
						"odinRendererAddInteractions",
						state.renderHandle,
						"Re-activate Airship Signaling Mast",
						"Re-activate Airship Signaling Mast",
						"", --"Re-activate Airship Signaling Mast",
						"", --"Re-activate Airship Signaling Mast",
						"",
						"",
						"click01",
						false,true)
				end
				
				send("rendInteractiveObjectClassHandler",
					"odinRendererAddInteractions",
					state.renderHandle,
					"Dismantle Airship Mast",
					"Dismantle Objects",
					"", --"Dismantle Airship Masts",
					"", --"Dismantle Objects",
					"hammer_icon",
					"construction",
					"Hammer Wood E",
					false,true)
			end
		end
		
		function airship_mast_reset_tooltip()
			local tooltipTitle = "Airship Mast"
			local tooltipDescription = "This airship mast is under construction."
		
			if not SELF.tags["under_construction"] then
				if SELF.tags["active_airship_mast"] then
					tooltipDescription = "This airship mast will signal to passing airships where they ought to drop off goods & immigrants."
				else -- if not active_airship_mast
					tooltipDescription = "This airship mast is currently disabled. If activated, goods and immigrants will be dropped nearby."
				end
			else
				
			end
			
			send("rendInteractiveObjectClassHandler",
                    "odinRendererBindTooltip",
                    state.renderHandle,
                    "ui//tooltips//groundItemTooltipDetailed.xml",
                    tooltipTitle,
                    tooltipDescription)
		end
	>>

	state
	<<
		int renderHandle
          int health
		int timer
		bool addedJob
          gameSimAssignmentHandle assignment
	>>

	receive Create(stringstringMapHandle init)
	<<
		printl("buildings", "airship_mast got Create")
		state.position.x = -1
		state.position.y = -1
		SELF.tags = { "under_construction" }
          state.health = 100
          state.addedJob = false
		state.assignment = nil
		state.timer = 100
		ready()
		sleep()
	>>

	receive GameObjectPlace( int x, int y ) 
	<<
		printl("buildings", "airship_mast got GameObjectPlace; " .. x .. " / " .. y )
		state.renderHandle = SELF.id
    
		local occupancyMap = 
			"..@..\\".. 
			".---.\\".. 
			"@-c-@\\"..
			".---.\\"..
			"..@..\\" 
		local occupancyMapRotate45 = 
			"..@..\\".. 
			".---.\\".. 
			"@-c-@\\"..
			".---.\\"..
			"..@..\\"
			
		send("gameSpatialDictionary", "registerSpatialMapString", SELF, occupancyMap, occupancyMapRotate45, false )
		send("gameSpatialDictionary", "gridAddObjectTo", SELF, state.position)
		send("gameSpatialDictionary", "toggleBuildingOnSquares", state.position.x, state.position.y, 1, 1, true)
	>>

	receive CompleteConstruction()
	<<
		printl("buildings", "airshipmast got CompleteConstruction")
		
		local tickerstring = "An airship signalling mast has been completed!"
		
		send("rendCommandManager",
			"odinRendererTickerMessage",
			tickerstring,
			"icon_airshiptower1",
			"ui\\orderIcons.xml")
		
		--[[send("rendCommandManager",
			"odinRendererStubMessage", --"odinRendererStubMessage",
			"ui\\orderIcons.xml", -- iconskin
			"icon_airshiptower1", -- icon
			"Airship Mast Constructed", -- header text
			tickerstring, -- text description
			"Right-click to dismiss.", -- action string
			"airshipMast", -- alert type (for stacking)
			"ui\\eventart\\cult_ritual.png", -- imagename for bg
			"low", -- importance: low / high / critical
			nil, -- object ID
			60 * 1000, -- duration in ms
			0, -- snooze
			nil)]]
		
		send("rendCommandManager",
			"odinRendererStubMessage",
			"ui\\orderIcons.xml", -- iconskin
			"icon_airshiptower1", -- icon
			"Airship Mast built", -- header text
			"Airship signalling mast completed! Airdrops will be sent to this location.", -- text description
			"Right-click to dismiss.", -- action string
			"airshipMastComplete", -- alert type (for stacking)
			"", -- imagename for bg
			"low", -- importance: low / high / critical
			state.renderHandle, -- object ID
			30 * 1000, -- duration in ms
			0, -- "snooze" time if triggered multiple times in rapid succession
			nil) -- gameobjecthandle of director, null if none

		local collection = query("gameObjectManager", "gameObjectCollectionRequest", "airdropMasts")[1]
		for k, otherMast in pairs(collection) do
			-- tell 'em to shut off. because swag
			send(otherMast, "turnOff", "another airship mast was built")
		end
		
		send("gameSession", "setSessionBool", "airdropOverride", true)
		send("gameSession", "setSessionBool", "airshipMastBuilt", true)
		send("gameSession", "setSessionInt", "airdropX", state.position.x)
		send("gameSession", "setSessionInt", "airdropY", state.position.y)
		
		collection[ #collection +1 ] = SELF
		
          SELF.tags = {"airship_mast", "destructible_wall", "active_airship_mast",}
		
          state.health = 40

		send("rendStaticPropClassHandler",
			"odinRendererCreateStaticPropRequest",
			SELF,
			"models/constructions/airshipMooringTower.upm",
			state.position.x,
			state.position.y)

		--[[send("rendStaticPropClassHandler",
			"odinRendererRotateStaticProp",
			state.renderHandle, rand(1,359), 0.25)--]]
          
          send("rendCommandManager",
               "odinRendererCreateParticleSystemMessage",
               "CeramicsStamperPoof",
               state.position.x,
               state.position.y )

		local occupancyMap = 
			"..@..\\".. 
			".ppp.\\".. 
			"@p#p@\\"..
			".ppp.\\"..
			"..@..\\"
			
		local occupancyMapRotate45 = 
			"..@..\\".. 
			".ppp.\\".. 
			"@p#p@\\"..
			".ppp.\\"..
			"..@..\\"
			
		send( "gameSpatialDictionary", "registerSpatialMapString", SELF, occupancyMap, occupancyMapRotate45, false )
          
		airship_mast_reset_controls()
		airship_mast_reset_tooltip()
	>>
	
	receive turnOff( string reason)
	<<
		if not SELF.tags["under_construction"] then
			if SELF.tags["active_airship_mast"] then
				send("gameSession", "setSessionBool", "airdropOverride", false)
				
				SELF.tags["active_airship_mast"] = nil
				
				send("rendCommandManager",
					"odinRendererTickerMessage",
					"Airship signalling mast disabled because " .. reason .. ". Airdrops will be deployed to a roughly central location.",
					"icon_airshiptower1",
					"ui\\orderIcons.xml")
				
				airship_mast_reset_controls()
				airship_mast_reset_tooltip()
			end
		end
	>>
	
	receive turnOn()
	<<
		if not SELF.tags["under_construction"] then
			if not SELF.tags["active_airship_mast"] then
				
				send("rendCommandManager",
					"odinRendererTickerMessage",
					"Airship signalling mast re-activated, airdrops will be sent to this position.",
					"icon_airshiptower1",
					"ui\\orderIcons.xml")
				
				SELF.tags["active_airship_mast"] = true
				send("gameSession", "setSessionBool", "airdropOverride", true)
				send("gameSession", "setSessionInt", "airdropX", state.position.x)
				send("gameSession", "setSessionInt", "airdropY", state.position.y)
				
				local collection = query("gameObjectManager", "gameObjectCollectionRequest", "airdropMasts")[1]
				for k, otherMast in pairs(collection) do
					if otherMast ~= SELF then
						send(otherMast, "turnOff", "another airship mast was enabled") 
					end
				end
				
				airship_mast_reset_controls()
				airship_mast_reset_tooltip()
			end
		end
	>>
	
	receive InteractiveMessage( string messagereceived )
	<<
		printl("buildings", "airship mast got InteractiveMessage: " .. messagereceived )
		if not state.addedJob then
			if messagereceived == "Dismantle Objects" then
				
				send(SELF,"turnOff", "this airship mast was ordered to be dismantled")
				
                    send("rendCommandManager",
                         "odinRendererCreateParticleSystemMessage",
                         "Small Beacon",
                         state.position.x,
                         state.position.y)
                    
				 results = query("gameBlackboard",
								"gameObjectNewAssignmentMessage",
								SELF,
								"Dismantle Objects",
								"construction",
								"construction")
				
                    state.assignment = results[1]
                    state.addedJob = true
                    send( "gameBlackboard",
                             "gameObjectNewJobToAssignment",
                             state.assignment,
                             SELF,
                             "Dismantle Objects",
                             "object",
                             true )
                         
				send("rendStaticPropClassHandler",
				    "odinRendererStaticPropExpressionMessage",
				    state.renderHandle,
				    "machine_thought64",
				    "jobaxe", false)
                    
                    send("rendInteractiveObjectClassHandler",
                              "odinRendererClearInteractions",
                              state.renderHandle)
               end
          end
	
		if messagereceived == "De-activate Airship Signaling" then
			send(SELF,"turnOff", "this mast was deactivated")
		elseif messagereceived == "Re-activate Airship Signaling Mast" then
			send(SELF,"turnOn")
		end
	>>
	
     receive InteractiveMessageWithAssignment( string messagereceived, gameSimAssignmentHandle assignment )
     <<
		printl("buildings", "airship mast got InteractiveMessageWithAssignment: " .. messagereceived )
          if not state.addedJob then
			if messagereceived == "Dismantle Objects" then
				
				send(SELF,"turnOff", "this airship mast was ordered to be dismantled")
				
                    send("rendCommandManager",
                         "odinRendererCreateParticleSystemMessage",
                         "Small Beacon",
                         state.position.x,
                         state.position.y)
                    
                    state.assignment = assignment
                    state.addedJob = true
                    send( "gameBlackboard",
                             "gameObjectNewJobToAssignment",
                             state.assignment,
                             SELF,
                             "Dismantle Objects",
                             "object",
                             true )
                         
				send("rendStaticPropClassHandler",
				    "odinRendererStaticPropExpressionMessage",
				    state.renderHandle,
				    "machine_thought64",
				    "jobaxe", false)
                    
                    send("rendInteractiveObjectClassHandler",
                              "odinRendererClearInteractions",
                              state.renderHandle)
               end
          end
		
		if messagereceived == "Cancel Construction" then
			send(SELF, "Clear", nil)
		end
		
		if messagereceived == "De-activate Airship Signaling" then
			send(SELF,"turnOff", "this airship mast was ordered to be deactivated")
		elseif messagereceived == "Re-activate Airship Signaling Mast" then
			send(SELF,"turnOn")
		end
     >>
     
	receive Update()
	<<

	>>

	respond isObstruction()
	<<
		return "obstructionResult", true
	>>

	receive JobCancelledMessage(gameSimJobInstanceHandle job)
	<<
		state.assignment = nil
		state.addedJob = false;
	>>

     receive ClearViolently( gameSimJobInstanceHandle ji, gameObjectHandle damagingObject )
     <<
          -- thunder & noise
		  --[[
          send("rendCommandManager",
                    "odinRendererCreateParticleSystemMessage",
                    "DustPuffXtraLarge",
                    state.position.x,
                    state.position.y )
			]]--
				
          send("rendInteractiveObjectClassHandler",
                    "odinRendererPlaySFXOnInteractive",
                    state.renderHandle,
                    "Break (Dredmor)" )
		
          local handle = query( "scriptManager",
                              "scriptCreateGameObjectRequest",
                              "objectcluster",
                              {legacyString = "Wrecked Airship Mast Cluster"} )[1]
		
		send(handle, "GameObjectPlace", state.position.x, state.position.y )
          
          send(SELF,"Clear", ji)
     >>
     
     receive Dismantle( gameSimJobInstanceHandle ji )
     <<
          -- spawn material, then delete self.
          local results = query("scriptManager",
							"scriptCreateGameObjectRequest",
							"item",
                                   {legacyString = "bricks"} )[1]
		
          send( results, "GameObjectPlace", state.position.x, state.position.y  )

          send("rendCommandManager",
                    "odinRendererCreateParticleSystemMessage",
                    "CeramicsStamperPoof",
                    state.position.x,
                    state.position.y )
               
          send(SELF,"Clear", ji)
     >>
     
	receive Clear( gameSimJobInstanceHandle ji )
	<<
		if SELF.tags["active_airship_mast"] then
			send("gameSession", "setSessionBool", "airdropOverride", false)
		else
			-- not the active mast, so whatevs.
		end
		
		local collection = query("gameObjectManager", "gameObjectCollectionRequest", "airdropMasts")[1]
		for k, otherMast in pairs(collection) do
			if otherMast == SELF then
				collection[k] = nil
				break
			end
		end
		
		-- are all masts destroyed?
		local num = 0
		for k, otherMast in pairs(collection) do
			num = num +1 	
		end
		if num == 0 then
			send("gameSession", "setSessionBool", "airshipMastBuilt", false)
		end

		send("gameSpatialDictionary", "toggleBuildingOnSquares", state.position.x, state.position.y, 1, 1, false);
          send("gameBlackboard", "gameObjectRemoveTargetingJobs", SELF, ji)
		send("rendStaticPropClassHandler", "odinRendererDeleteStaticProp", SELF.id)
		send("gameSpatialDictionary", "gridRemoveObject", SELF)

		send("rendCommandManager",
			"odinRendererCreateParticleSystemMessage",
			"DustPuffMassive",
			state.position.x,
			state.position.y)
				
		destroyfromjob(SELF,ji)
	>>

	receive JobCancelledMessage( gameSimJobInstanceHandle ji )
	<<
		if SELF.tags["under_construction"] then
			send(SELF,"Clear", ji)
		else
			send("rendInteractiveObjectClassHandler",
                         "odinRendererClearInteractions",
                         state.renderHandle)
			
			send("rendInteractiveObjectClassHandler",
                    "odinRendererAddInteractions",
                    state.renderHandle,
                    "Dismantle Object",
                    "Dismantle Objects",
                         "", --"Dismantle Object",
                         "", --"Dismantle Objects",
                         "hammer_icon",
                                 "construction",
                                  "Hammer Wood E",
						    false,
						    true)
		end
	>>

>>