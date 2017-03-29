gameobject "graveyard"
<<
	local 
	<<
		
	>>

	state
	<<
		table contents
		table empty_squares
		table tilled_squares
		int x
		int y
		int w
		int h
		bool addedJob
		gameSimAssignmentHandle assignment
		gameGridPosition position
	>>

	receive Create( stringstringMapHandle init )
	<<
		--printl("creating a graveyard.")
		state.assignment = nil
		ready()
	>>

	receive Update()
	<<

	>>

	receive odinCreateZoneMessage(string myName, string createString,  int x, int y, int w, int h, int originX, int originY, bool fillSquares)
	<<
		printl("buildings", "creating a graveyard... part II")
		send("gameZoneManager", "ZoneNewZoneMessage", SELF, "graveyard", "Graveyard");

		state.x = x
		state.y = y
		state.w = w
		state.h = h
		state.position.x = state.x
		state.position.y = state.y
		
		state.empty_squares = {}

		for i = 1,w-1 do
			for j = 1,h-1 do
				if (i-1)%2 == 1 and (j-1)%2 == 1 then
					gp = gameGridPosition:new()
					gp.x = state.x + i
					gp.y = state.y + j
					state.empty_squares[#state.empty_squares+1] = gp
				end
			end
		end

		send("rendOdinZoneClassHandler", "odinRendererCreateZoneRequest", SELF, state.x,state.y,state.x+state.w,state.y+state.h, "Graveyard");

		-- tell an overseer to show up and get this thing running
		send("rendInteractiveObjectClassHandler",
				"odinRendererBindTooltip",
				SELF.id,
				"ui//tooltips//groundItemTooltipDetailed.xml",
				"Graveyard",
				"Your fallen characters will be taken here when the others have a minute." )
	>>

	respond gridReportPosition() 
	<<
		return "gridReportedPosition", state.position
	>>

	respond gridGetPosition()
	<<
		return "reportedPosition", state.position
	>>

	respond ZoneGetTargetSquare()
	<<
		if (#state.empty_squares < 1) then
			return "ZoneTargetSquare", -1, -1
		else
			return "ZoneTargetSquare", state.empty_squares[1].x, state.empty_squares[1].y
		end
	>>

	respond ZoneGetParams ()
	<<
		return "ZoneParams", state.x, state.y, state.w, state.h, SELF
	>>

	receive ZoneRemoveSquare ( int x, int y )
	<<
		--printl("buildings", "graveyard trying to remove square " .. tostring(x) .. "," .. tostring(y))
		for i = 1,#state.empty_squares do
			if (state.empty_squares[i].x == x and state.empty_squares[i].y == y) then	
				table.remove(state.empty_squares,i)
				--printl("buildings"," graveyard removnig square " .. tostring(x) .. "," .. tostring(y))
				break		
			end
		end
	>>

	receive SquareAvailable ( int x, int y )
	<<
		-- printl("FARM: new square available.");
		gp = gameGridPosition:new()
		gp.x = x
		gp.y = y
		state.empty_squares[#state.empty_squares+1] = gp
	>>
>>