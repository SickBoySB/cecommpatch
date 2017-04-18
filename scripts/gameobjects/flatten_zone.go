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
		
		-- CECOMMPATCH 
		-- this block fixes crashes from flattening terrain on the edge of the map by preventing flattening jobs from occurring there
		-- some intelligent resizing of the job occurs so that only valid areas are selected (rather than just refusing the job entirely)
		
		local l_min = 5
		local l_max = 250
		
		local newsx = state.x
		local newsy = state.y
		local newx = x
		local newy = y
		local neww = w
		local newh = h

		-- change x if needed
		if state.x < l_min then
			newsx = l_min
			newx = l_min
		elseif state.x > l_max then
			newsx = l_max
			newx = l_max
		end
		
		neww = w - (math.abs(state.x - newsx))
		
		if newsx + neww > l_max then
			neww = l_max - newsx
		end
		
		-- now do the same for y
		if state.y < l_min then
			newsy = l_min
			newy = l_min
		elseif state.y > l_max then
			newsy = l_max
			newy = l_max
		end
		
		newh = h - (math.abs(state.y - newsy))
		
		if newsy + newh > l_max then
			newh = l_max - newsy
		end
		
		
		if (newsx > state.x + w) or 
			(newsx + neww > state.x + w) or
			(neww < 1) then
			neww = 0 -- no valid spots
		end
		
		if (newsy > state.y + h) or 
			(newsy + newh > state.y + h) or
			(newh < 1) then
			 newh = 0 -- no valid spots
		end
		
		state.x = newsx
		state.y = newsy
		x = newx
		y = newy
		w = neww
		h = newh
		-- /CECOMMPATCH FIX
		if neww > 0 and newh > 0 then
		
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