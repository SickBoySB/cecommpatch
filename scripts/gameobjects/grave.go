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
		
		local models = {
			"graveyardHeadstone00.upm",
			"graveyardHeadstone01.upm",
			"graveyardHeadstone02.upm",
			"graveyardCogStone00.upm",
			"graveyardCogStone01.upm",
			"graveyardCogWood00.upm",
			--"graveyardObelisk00.upm",
			--"graveyardObelisk01.upm",
		}
		send("rendStaticPropClassHandler",
			"odinRendererCreateStaticPropRequest",
			SELF,
			"models/constructions/" .. models[rand(1,#models)],
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