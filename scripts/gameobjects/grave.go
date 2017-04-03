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
		local models = {
			-- vanilla headstones
			"graveyardHeadstone00.upm",
			"graveyardHeadstone01.upm",
			"graveyardHeadstone02.upm",
			"graveyardCogStone00.upm",
			"graveyardCogStone01.upm",
			"graveyardCogWood00.upm",
			-- CECOMMPATCH additions.. mostly unused fence posts
			"fences/basicRailFencePost01.upm",
			"fences/basicRailFencePost02.upm",
			"fences/basicRailFencePost03.upm",
			"fences/fencePost01.upm",
			"fences/fencePost02.upm",
			"fences/rusticFencePost01.upm",
			"fences/rusticFencePost02.upm",
			"fences/rusticFencePost03.upm",
			"fences/whitePicketFencePost01.upm",
			"memorial00.upm",
			-- BIG graves - disabled for now. they have a weird alignment
			--"graveyardObelisk00.upm",
			--"graveyardObelisk01.upm",
		}
		
		-- TODO: base randomization on class/faction
		local grave_model = models[rand(1,#models)]
		local grave_rotate = rand(-7,7) -- slight rotation for variety. too much is weird looking
		
		-- memorial00.upm has a bad rotation by default, hackish fix
		if grave_model == "memorial00.upm" then
			grave_rotate = grave_rotate + 180
		end
		
		state.position.x = x
		state.position.y = y
		state.renderHandle = SELF.id
		
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
          local tooltipDescription = "A marker to show where a rich vein of bones and rotting meat may be found."
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