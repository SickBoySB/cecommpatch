event "foreign_intervention"
<<
	local
	<<
		function choose_faction()
			local o = ""
			local dominant = query("gameSession","getSessionString", "dominantFaction")[1]
			
			-- CECOMMPATCH - hack to fix dominantFaction pulling fullName rather than shortName
			-- TODO: do it properly via EntityDB lookup
			
			local dominant_short_name = {
				["Grossherzoginnentum von Stahlmark"] = "Stahlmark",
				["Novorus Imperiya"] = "Novorus",
				["Republique Mecanique"] = "Republique",
				-- not sure if any of the below is ever used, but just to be safe..
				["The Clockwork Empire"] = "Empire",
				["Fishpeople"] = "Fishpeople",
				["Bandits"] = "Bandits"
			}		
			
			local nations = { "Stahlmark", "Novorus", "Republique" }
			
			if dominant_short_name[dominant] then
				table.insert(nations,dominant_short_name[dominant])
			end
			
			o = nations[rand(1,#nations)]
			
			return o
		end
	>>
	
	state 
	<<
		bool alertTriggered
          int counter
          table spawnLocation
		string dialogBoxResults
		int timeout
	>>

	receive Create( stringstringMapHandle name )
	<< 
		printl("events","Foreign intervention event started")
		
		-- choose a nation
		-- have them send a squad to do a mission based on relation status.
		
		state.nationName = choose_faction()
		state.nation = query("gameSession","getSessiongOH", state.nationName)[1]
		state.nationInfo = EntityDB[ state.nationName .. "Info"]
	>>

	receive boxSelectionResult( string id, string message )
	<<
		state.dialogBoxResults = message
	>>

	receive alertSelectedResult( string message )
	<<
		state.alertTriggered = true
	>>
	
	FSM 
	<<
		["start"] = function(state,tags)
			if not state.nation then
				-- OLDER SAVE FIX
				-- pick a new faction since dominantFaction was broken previously. shouldn't ever hit this with new saves
				printl("CECOMMPATCH: foreign_intervention event broke.. invalid faction. retrying")
				state.nationName = choose_faction()
				state.nation = query("gameSession","getSessiongOH", state.nationName)[1]
				state.nationInfo = EntityDB[ state.nationName .. "Info"]
				-- /OLDER SAVE FIX
				
				return "start"
			end	
			
			settimer("Foreign Troops Event Timer", 0)

			local isPatrol = false
			if rand(1,2) == 1 then isPatrol = true end
			
			local isNeutral = true
			local isAllied = query(state.nation,"isFriendly")[1] 
			local isHostile = query(state.nation,"isHostile")[1]
			if isAllied or isHostile then
				isNeutral = false
			end
			
			printl("events", "starting foreign_intervention. Nation= " .. state.nationName .. ", allied: " .. tostring(isAllied) ..
				  " / hostile: " .. tostring(isHostile) .. " / neutral: " .. tostring(isNeutral) )
			
			local numGroup = "NUMBER_GROUP" -- query(group,"getNumMembers")[1] -- sub this in post-spawn.
			
               state.alertTriggered = false
               local mission = ""
			local s = state.nationInfo.adjective .. " troops is" 
			if isAllied then
				s = "A unit of ".. numGroup .." allied " .. s
				if isPatrol then
					s = s .. " passing through the area."
					mission = "patrol"
				else
					s = s .. " going to visit our settlement."
					mission = "visit"
				end
			elseif isHostile then
				s = "A unit of ".. numGroup .." hostile " .. s
				if isPatrol then
					s = s .. " passing through the area - watch out!"
					mission = "patrol"
				else
					s = s .. " moving to attack our colony!"
					mission = "attack"
				end
			else
				s = "A unit of ".. numGroup .." neutral " .. s
				if isPatrol then
					s = s .. " passing through the area."
					mission = "patrol"
				else
					s = s .. " scouting around our colony. Stand vigilant."
					mission = "scout"
				end
			end
			
			function findMapEdgeSpawnLoc()
				local valid = false
				local i = 1
				
				local x_max = query("gameSession","getSessionInt","x_max")[1]
				local y_max = query("gameSession","getSessionInt","y_max")[1]
		
				local spawnRects = {
						[1] = {x=8,y=8,w=x_max-16,h=16, name = "West"},
						[2] = {x=8,y=8,w=16,h=y_max-16, name = "Northwest"},
						[3] = {x=8,y=y_max-16,w=x_max-16,h=16, name = "Northeast"},
						[4] = {x=x_max-16,y=8,w=16,h=y_max-16, name = "East"},
						}
					
				local activeRect = spawnRects[ rand(1,#spawnRects) ]
				local newLoc = gameGridPosition:new()
					
				while not valid do
					newLoc.x = rand(activeRect.x, activeRect.x + activeRect.w )
					newLoc.y = rand(activeRect.y, activeRect.y + activeRect.h )
					
					local iswater = query("gameSpatialDictionary",
									  "gridHasSpatialTag",
									  newLoc,
									  "water" )[1]
					
					if not iswater then
						-- yes!
						valid = true
						printl("events",
							  "bandit_spawn : found valid location, returning " ..
							  tostring(newLoc.x) .. ", " .. tostring(newLoc.y) )
						
						return newLoc
					elseif i == 150 then
						-- end it.
						printl("events", " couldn't find spawn position!")
						
						valid = true
						return false
					end
					i = i + 1
				end
			end
			
			-- place foreign unit

			local startLoc = findMapEdgeSpawnLoc()
			if not startLoc then
				return "abort",true
			end
			
               local foreign_group = query( "scriptManager",
									"scriptCreateGameObjectRequest",
									"foreigner_group",
									{ legacyString = state.nationName .." Unit",
									mission = "idle", } )[1]
			
               send(foreign_group,
				"GameObjectPlace",
				startLoc.x,
				startLoc.y )

			local idleTime = 30
			if mission == "visit" then idleTime = 60 end
			if mission == "scout" then idleTime = 30 end
			
               send(foreign_group,"pushMission", mission, nil, idleTime)

			-- Ugh, messy.
			s = string.gsub(s, "NUMBER_GROUP", query(foreign_group,"getNumMembers")[1] )
			
			-- Okay, NOW tell them we did it, after successful placement
			send("rendCommandManager",
				"odinRendererTickerMessage",
				s,
                    state.nationInfo.iconTroops,
                    "ui\\thoughtIcons.xml")
			
			if isAllied then
				
				local leader = query(foreign_group,"getLeader")[1]
				local leaderhandle = query(leader,"ROHQueryRequest")[1]
				
				send("rendCommandManager",
					"odinRendererStubMessage",
					"ui\\thoughtIcons.xml", -- iconskin
					state.nationInfo.iconTroops, -- icon
					"Allies Arrive", -- header text
					s, -- text description
					"Left-click to zoom. Right-click to dismiss.", -- action string
					"foreignIntervention", -- alert type (for stacking)
					"", -- imagename for bg
					"low", --"high", -- importance: low / high / critical
					leaderhandle, -- object ID
					60 * 1000, -- duration in ms
					0, -- "snooze" time if triggered multiple times in rapid succession
					nil)
				
				send("rendCommandManager",
					"odinRendererPlaySoundMessage",
					"alertNeutral")
				
			elseif isHostile then
				send("rendCommandManager",
					"odinRendererStubMessage",
					"ui\\thoughtIcons.xml", -- iconskin
					state.nationInfo.iconTroops, -- icon
					"Hostile Foreigners!", -- header text
					s, -- text description
					"Right-click to dismiss.", -- action string
					"foreignIntervention", -- alert type (for stacking)
					"", -- imagename for bg
					"low", --"high", -- importance: low / high / critical
					nil, -- leaderhandle, -- object ID
					60 * 1000, -- duration in ms
					0, -- "snooze" time if triggered multiple times in rapid succession
					nil)
				
				send("rendCommandManager",
					"odinRendererPlaySoundMessage",
					"alertDanger")
				
				if mission == "attack" or mission == "patrol" then
					
					local daycount = query("gameSession","getSessionInt","dayCount")[1]
					local disturbance = query("gameSession","getSessionInt","disturbancePoints")[1]
					if daycount > 40 then
						if disturbance > 150 then
							-- you deserve some pain.
							 local foreign_group = query( "scriptManager",
									"scriptCreateGameObjectRequest",
									"foreigner_group",
									{ legacyString = state.nationName .." Unit",
									mission = "idle", } )[1]
			
							send(foreign_group, "GameObjectPlace", -1, -1 )
							send(foreign_group,"pushMission", mission, nil, 60)
					
							send("gameSession","incSessionInt","disturbancePoints",-5)
						end
					end
				end
			else
				send("rendCommandManager",
					"odinRendererStubMessage",
					"ui\\thoughtIcons.xml", -- iconskin
					state.nationInfo.iconTroops, -- icon
					"Neutral Foreigners", -- header text
					s, -- text description
					"Right-click to dismiss.", -- action string
					"foreignIntervention", -- alert type (for stacking)
					"", -- imagename for bg
					"low", --"high", -- importance: low / high / critical
					nil, -- leaderhandle, -- object ID
					60 * 1000, -- duration in ms
					0, -- "snooze" time if triggered multiple times in rapid succession
					nil)
				
				send("rendCommandManager",
					"odinRendererPlaySoundMessage",
					"alertNeutral")
			end
			
			return "final", true
		end,

		["final"] = function(state,tags)  
			return
		end,

		["abort"] = function(state,tags)  
			return
		end
	>>
>>
