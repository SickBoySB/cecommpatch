gameobject "grave"
<<
	local 
	<<

	>>

	state
	<<
		gameGridPosition position
		int renderHandle
		string gabionType
		bool addedJob
		int corpseID
	>>

	receive Create(stringstringMapHandle init)
	<<
		local type = init["legacyString"]

		state.position.x = -1
		state.position.y = -1
		state.gabionType = type
		SELF.tags = {"grave"}

		ready()
		sleep()
	>>

	receive GameObjectPlace( int x, int y ) 
	<<	
		state.position.x = x
		state.position.y = y
		state.renderHandle = SELF.id
		
		local models_lc = {
			"graveyardCogWood00.upm",
			"graveyardCogStone00.upm",
			"graveyardCogStone01.upm",
			"fences/fencePost01.upm",
			--"fences/fencePost02.upm",
		}
		
		local models_mc = {
			"graveyardHeadstone00.upm",
			"graveyardHeadstone01.upm",
			"graveyardHeadstone02.upm",
		}
		
		local models_uc = {
			"graveyardObelisk00.upm",
			"graveyardObelisk01.upm",
		}
		
		local models_military_lc = {
			"fences/whitePicketFencePost01.upm",
		}
		
		local models_military_mc = {
			"memorial00.upm",
		}
		
		local models_other = {
			"fences/rusticFencePost01.upm",
			"fences/rusticFencePost02.upm",
			"fences/rusticFencePost03.upm",
			"fences/basicRailFencePost01.upm",
			"fences/basicRailFencePost02.upm",
			"fences/basicRailFencePost03.upm",
		}
			
		local models = {}
		local grave_who = "Unknown"

		local results = query("gameSpatialDictionary",
						  "allObjectsInRadiusWithTagRequest",
						  state.position,
						  1,"corpse",true)[1]
		
		if results then
			grave_who_tags = query(results[1],"getTags")[1]
			grave_who = query(results[1],"getName")[1]
			--printl("CECOMMPATCH: " .. grave_who)
		end	
		
		if grave_who_tags["citizen"] then
			if grave_who_tags["lower_class"] then
				if grave_who_tags["military"] then
					models = models_military_lc
				else
					models = models_lc
				end
			elseif grave_who_tags["middle_class"] then
				if grave_who_tags["military"] then
					models = models_military_mc
				else
					models = models_mc
				end
			elseif grave_who_tags["upper_class"] then
				models = models_uc
			end
		else
			-- non-citizen, give 'em sticks
			models = models_other
		end
		
		-- TODO: base randomization on class/faction
		local grave_model = models[rand(1,#models)]
		local grave_rotate = rand(-7,7) -- slight rotation for variety. too much is weird looking
		
		-- memorial00.upm has a bad rotation by default, hackish fix
		if grave_model == "memorial00.upm" then
			grave_rotate = grave_rotate + 180
		end
		
		send("rendStaticPropClassHandler",
			"odinRendererCreateStaticPropRequest",
			SELF,
			"models/constructions/" .. grave_model,
			state.position.x,
			state.position.y)

		-- TODO: check with Micah to see if these are, in fact, sane (I would prefer borders.)

		--[[local occupancyMap = 
		"PPPPP\\".. 
		"PpCpP\\"..
		"PPPPP\\"

		local occupancyMapRotate45 = 
		"PPPPP\\".. 
		"PpCpP\\"..
		"PPPPP\\"]]
		
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
			
		
		send("gameSpatialDictionary",
               "registerSpatialMapString",
               SELF,
               occupancyMap,
               occupancyMapRotate45,
               true )
         
		send("gameSpatialDictionary",
               "gridAddObjectTo",
               SELF,
               state.position)
		
		send("rendCommandManager",
			"odinRendererHighlightSquare",
			state.position.x,
			state.position.y,
			255, 255, 255, true)
		
		send("rendStaticPropClassHandler",
			"odinRendererMoveStaticPropGridWithHeight",
			state.renderHandle,
			state.position,
			0)
		
		send("rendStaticPropClassHandler",
			"odinRendererOrientStaticProp",
			state.renderHandle,
			0)
		
		-- some slight rotation... must be done LAST because of odinRendererOrientStaticProp
		send("rendStaticPropClassHandler",
			"odinRendererRotateStaticProp",
			state.renderHandle,
			grave_rotate,
			0.25)
		
		-- TODO: get some interesting grave text, randomized, maybe based on the colonist being buried?
          local tooltipTitle = "Grave"
          local tooltipDescription = "Here rests " .. grave_who .. ". Good riddance."
          send("rendInteractiveObjectClassHandler",
                    "odinRendererBindTooltip",
                    state.renderHandle,
                    "ui//tooltips//groundItemTooltipDetailed.xml",
                    tooltipTitle,
                    tooltipDescription)

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

	respond gridGetPosition()
	<<
		return "reportedPosition", state.position
	>>

	respond gridReportPosition()
	<<
		return "gridReportedPosition", state.position
	>>
	
>>