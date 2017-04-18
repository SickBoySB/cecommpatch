gameobject "flatten_zone"
<<
	local 
	<<
		
	>>

	state
	<<
		table flattenSquares
		int x
		int y
		int w
		int h
		int renderHandle
		int targetHeight
		bool addedJob
		gameSimAssignmentHandle assignment
		gameGridPosition position
	>>

	receive Create( stringstringMapHandle init )
	<<
		state.assignment = nil
		ready()		
	>>

	receive Update()
	<<
	>>

	receive odinCreateZoneMessage(string myName, string createString, int x, int y, int w, int h, int originX, int originY, bool fillSquares)
	<<
		state.x = x
		state.y = y
		state.w = w
		state.h = h
		state.position.x = originX
		state.position.y = originY
		
		state.targetHeight = query("gameSpatialDictionary", "gridGetHeight", state.position)[1];

		state.flattenSquares = {};

		local results = query("gameBlackboard",
					 "gameObjectNewAssignmentMessage",
					 SELF,
					 "Flatten Terrain",
					 "flattening",
					 "construction")
		
		state.assignment = results[1]		
		state.renderHandle = SELF.id;
		send("gameBlackboard", "SetAssignmentFlatteningInformation", state.assignment, x, y, w, h, state.targetHeight);

		send("rendBeaconClassHandler",
			"CreateAssignmentBeacon",
			state.assignment,
			"flatten_icon",
			"ui\\thoughtIcons.xml",
			"ui\\thoughtIconsGray.xml",
			4)
				 
		for i = 0,w do
			for j = 0,h do
				gp = gameGridPosition:new()
				gp.x = state.x + i
				gp.y = state.y + j

				send("rendBeaconClassHandler",
					"AddPositionToAssignmentBeacon",
					state.assignment,
					gp.x,gp.y)
			end
		end
	>>

	receive InteractiveMessage( string messagereceived )
	<<
		printl ("Message Received: " .. messagereceived );

		if messagereceived == "Cancel Job" then
			send("gameBlackboard", "cancelAssignment", state.assignment)
		end

	>>
	respond gridReportPosition() 
	<<
		return "gridReportedPosition", state.position
	>>

	respond gridGetPosition()
	<<
		return "reportedPosition", state.position
	>>
	
	respond zoneCheckSpace(gameGridPosition pos)
	<<
		local height = query("gameSpatialDictionary", "gridGetHeight", pos)[1];
		
		return "zoneCheckSpaceResponse", height ~= state.targetHeight;
	>>

	respond ZoneGetTargetSquare()
	<<
		if (#state.flattenSquares < 1) then
			return "ZoneTargetSquare", -1, -1
		else
			return "ZoneTargetSquare", state.flattenSquares[1].x, state.flattenSquares[1].y
		end
	>>

	respond ZoneGetParams ()
	<<
		return "ZoneParams", state.x, state.y, state.w, state.h, SELF
	>>

	receive ZoneRemoveSquare ( int x, int y )
	<<
		printl("trying to remove square " .. tostring(x) .. "," .. tostring(y))
		for i = 1,#state.flattenSquares do
			if (state.flattenSquares[i].x == x and state.flattenSquares[i].y == y) then	
				table.remove(state.flattenSquares,i)
				printl("removing square " .. tostring(x) .. "," .. tostring(y))
				break		
			end
		end
	>>
	
	receive AssignmentSuspendedMessage(gameSimAssignmentHandle a)
	<<
		local done = true;
	
		gp = gameGridPosition:new()
		for i = 0,state.w do
			for j = 0,state.h do
				gp.x = state.x + i
				gp.y = state.y + j

				local height = query("gameSpatialDictionary", "gridGetHeight", gp)[1];
				if height ~= state.targetHeight then
					done = false;
				end
			end
		end
		
		if done then
			send("gameBlackboard", "cancelAssignment", a)
			
			-- delete self
			send("rendOdinZoneClassHandler", "odinRendererDeleteZoneRequest", SELF);
			send("gameSpatialDictionary", "gridRemoveObject", SELF);
			destroy(SELF);
		end
	>>

	receive AssignmentCancelledMessage(gameSimAssignmentHandle a)
	<<
		send("rendOdinZoneClassHandler", "odinRendererDeleteZoneRequest", SELF);
		send("gameSpatialDictionary", "gridRemoveObject", SELF);
	
		destroy(SELF);
	>>
	
	respond GetTargetHeight()
	<<
		return "TargetHeightResponse", state.targetHeight
	>>
	
	receive SquareAvailable ( int x, int y )
	<<
		gp = gameGridPosition:new()
		gp.x = x
		gp.y = y
		state.empty_squares[#state.flattenSquares+1] = gp
	>>
>>