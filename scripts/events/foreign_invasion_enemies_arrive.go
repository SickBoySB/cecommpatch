event "foreign_invasion_enemies_arrive"
<<
	state 
	<<
		string invader
	>>

	receive Create( stringstringMapHandle init )
	<<
		printl("events","foreign_invasion_enemies_arrive started")
		
		state.director = query("gameSession","getSessiongOH","event_director" .. init.director_name)[1]
		state.invader = query(state.director,"getKeyString","invader")[1]
		state.invaderFaction = query("gameSession","getSessiongOH", state.invader )[1]
		state.invaderNationInfo = EntityDB[ state.invader .. "Info"]
		state.counter = 0
		
		state.dialogBoxResults = ""
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
			-- Check to see if relations are still hostile. If not, end via peace made path.
			if not query(state.invaderFaction,"isHostile")[1] then
				printl("events", "foreign_invasion: wait a sec, " .. state.invader .. " is at peace with us! ")
				state.counter = 600
				return "suddenly_peace", true
			end
			
			state.counter = 0
			
			return "waiting1"
		end,
		
		["waiting1"] = function(state,tags)
			if state.counter >= 60 then
				printl("events", "foreign_invasion_enemies_arrive finished waiting1, doing enemyspawn")
				return "enemyspawn"
			end
			state.counter = state.counter + 1
			
			return "waiting1"
		end,

		["enemyspawn"] = function(state,tags)

			-- spawn enemies.
			-- Find a quiet place to spawn.
			function fakeSpawnNearColony()
				
				local x_max = query("gameSession","getSessionInt","x_max")[1]
				local y_max = query("gameSession","getSessionInt","y_max")[1]
		
				local rect = {x=8,y=8,w=x_max-8,h=y_max-8}
				local spawnValid = false
				local i = 0
				while not spawnValid do
					local newLoc = gameGridPosition:new()
					newLoc.x = rand(rect.x, rect.x + rect.w )
					newLoc.y = rand(rect.y, rect.y + rect.h )
				
					local civ = query( "gameSpatialDictionary", "gridGetCivilization", newLoc )[1]
					local iswater = query( "gameSpatialDictionary", "gridHasSpatialTag", newLoc, "water" )[1]
					
					local observerNearbyResults = query("gameSpatialDictionary",
									   "isObjectInRadiusWithTag",
									   newLoc,
									   15,
									   "observer")[1]
					
					local friendlyNearbyResults = query("gameSpatialDictionary",
									   "isObjectInRadiusWithTag",
									   newLoc,
									   15,
									   "friendly_agent")[1]
					
					local citizenNearbyResults = query("gameSpatialDictionary",
									   "isObjectInRadiusWithTag",
									   newLoc,
									   15,
									   "citizen")[1]
					
					if civ > 10 and
						not iswater and
						not observerNearbyResults and
						not friendlyNearbyResults and
						not citizenNearbyResults then
						
						spawnValid = true
						return newLoc
					end
					
					i = i + 1
					if i > 256 then
						return false
					end
				end
			end
			
			function findValidLocInRect( rect )
				local spawnValid = false
				local i = 0
				while not spawnValid do
					local newLoc = gameGridPosition:new()
					newLoc.x = rand(rect.x, rect.x + rect.w )
					newLoc.y = rand(rect.y, rect.y + rect.h )
				
					local civ = query( "gameSpatialDictionary", "gridGetCivilization", newLoc )[1]
					local iswater = query( "gameSpatialDictionary", "gridHasSpatialTag", newLoc, "water" )[1]
					
					if civ > 0 then
						-- no!
					elseif iswater then
						-- no!
					elseif i > 200 then
						-- wow, okay. Just give up.
						return false
					else	
						spawnValid = true
						printl("events",
							  "foreign_invasion : found valid location, returning " ..
							  tostring(newLoc.x) .. ", " .. tostring(newLoc.y) )
						
						return newLoc
					end
					i = i + 1
				end
			end
			
			
			local numTroops = 0
			local numUnits = rand(3,4)
			
			local dominant = query("gameSession","getSessionString", "dominantFaction")[1]
			if state.invader == dominant then
				numUnits = numUnits + rand(1,2)
			end
			
			send(state.director,"setKeyInt","enemy_squad_count", numUnits )
			
			--local startLoc = findValidLocInRect( state.spawnRects[startPos], false )
			
			local startLoc = fakeSpawnNearColony()
			
			local x_max = query("gameSession","getSessionInt","x_max")[1]
			local y_max = query("gameSession","getSessionInt","y_max")[1]
				
				
			local startX = query("gameSession", "getSessionInt", "startX")[1]
			local startY = query("gameSession", "getSessionInt", "startY")[1]
			
			local searchX = startX - math.floor(x_max * 0.25)
			local searchY = startY - math.floor(y_max * 0.25)
			
			if searchX < 2 then searchX = 2 end
			if searchY < 2 then searchY = 2 end
			
			local searchW = math.floor(x_max * 0.5) -16
			local searchH = math.floor(y_max * 0.5) -16
			
			if searchW + searchX > x_max then searchX = x_max - searchW end
			if searchH + searchY > y_max then searchY = y_max - searchH end
				
			printl("events", " search x/y/w/h = " .. searchX .. " / " .. searchY .. " / " .. searchW .. " / " .. searchH )
				
			local endLoc = findValidLocInRect( {	x= searchX,
											y= searchY,
											w= searchW,
											h= searchH},
											true )
			local mission = "attack"
			
			if startLoc == false or endLoc == false then
				-- um, failed to place something. Abort!
				printl("events", "WARNING: aborting foreign invasion due to failure to place start/end!")
				
				local s = "The " .. state.invaderNationInfo.adjective .. " troops got lost on the way to our settlement. \z
					It looks like there won't be an invasion after all. Huzzah!"
					
				send(state.director,
					"forwardEventStatusText",
					s,
					s)
				
				send(state.director, "setKeyString", "ending", "abort")
				return "abort"
			end
			
			
			for i=1,numUnits do

				local foreign_group = query( "scriptManager",
										"scriptCreateGameObjectRequest",
										"foreigner_group",
										{ 	legacyString = state.invader .." Unit",
											nationName = state.invader,
											invasion_director = "mundane", 
											mission = "idle", } )[1]
				
				send(foreign_group,"GameObjectPlace", startLoc.x + rand(-8,8), startLoc.y + rand(-8,8) )
				send(foreign_group,"pushMission", mission, endLoc, 300)
				
				numTroops = numTroops + query(foreign_group,"getNumMembers")[1]
			end
			
			-- Okay, NOW tell them enemies are arriving we did it, after successful placement
			local dialog_message = "The " .. state.invaderNationInfo.adjective .. " soldiers are advancing, prepare the defenses! \n\n\z
				We expect to face " .. numUnits .. " squads totalling " .. numTroops .." enemy troops.\n\n\z
				The Empire expects that each of our brave soldiers will do their duty."
			
			send(state.director,"setKeyInt","numTroops",numTroops)
			
			send(state.director,
				"forwardEventStatusText",
				"The " .. state.invaderNationInfo.adjective .. " arrived to attack our colony.",
				"The " .. state.invaderNationInfo.adjective .. " arrived to attack our colony.")
			
			send("scriptUIManager",
				"createSelectionDialogBox",
				SELF, 
				"foreignInvastionResponse", 
				"Foreign Invasion!", 
				dialog_message, 
				{
					{ 	["buttonText"] = "To arms! We shall win the day!",
						["toolTipText"] = "I hope this works.",
						["dialogBoxResults"] = "",
					},
				},
				"ui\\eventart\\battle.png")
			
			send("rendCommandManager",
                         "odinRendererPlaySoundMessage",
                         "alertDanger")
			
			send("rendCommandManager",
				"odinRendererStubMessage",
				"ui\\thoughtIcons.xml", -- iconskin
				state.invaderNationInfo.iconTroops, -- icon
				"Enemies Arrive", -- header text
				"The " .. state.invaderNationInfo.adjective .. " soldiers are advancing on our colony!", -- text description
				"Right-click to dismiss.", -- action string
				"invasionArrivesStub", -- alert type (for stacking)
				"", -- imagename for bg
				"low", -- importance: low / high / critical
				nil, -- object ID
				60 * 1000, -- duration in ms
				0,state.director) -- snooze
			
			state.counter = 300
			return "invasion_choice2"
		end,

		["invasion_choice2"] = function(state,tags)
			state.counter = state.counter - 1
			if state.counter == 0 then
				-- make a choice.
				local choicetable = {
					{ 	["buttonText"] = "Request a bombing run from the Imperial Air Corps.",
						["toolTipText"] = "Death from above! That's what I always say.",
						["dialogBoxResults"] = "stage2_bombing",
					},
					{ 	["buttonText"] = "Have the War Office send us additional troops.",
						["toolTipText"] = "Our brave redcoats will help hold the lines.",
						["dialogBoxResults"] = "stage2_reinforcements",
					},
				}

				send("scriptUIManager",
					"createSelectionDialogBox",
					SELF, 
					"foreignInvastionResponse", 
					"Assistance from the War Office", 
					"With the foreign invasion under way, we can jolly well request additional assistance from the War Office. \n\n\z
					Shall we have the Imperial Air Corps bomb the enemy or request that they drop off a squad of Her Majesty's Finest Redcoats?", 
					choicetable,
					"ui\\eventart\\airship_flying_away.png")
				
				state.counter = 0
				return "waiting"
			end
			
			return "invasion_choice2"
		end,

		["waiting"] = function(state,tags)
		
			if state.dialogBoxResults ~= "" and
				state.dialogBoxResults ~= nil then
				
				if state.dialogBoxResults == "redcoats" then
					
					send("rendCommandManager",
						"odinRendererStubMessage",
						"ui\\orderIcons.xml", -- iconskin
						"icon_rally", -- icon
						"Redcoats Requested", -- header text
						"We have requested a squad of trained Redcoats from The War Office.", -- text description
						"Right-click to dismiss.", -- action string
						"invasionArrivesResponseStub", -- alert type (for stacking)
						"", -- imagename for bg
						"low", -- importance: low / high / critical
						nil, -- object ID
						60 * 1000, -- duration in ms
						0, state.director) -- snooze
					
				elseif state.dialogBoxResults == "bombing" then
					
					send("rendCommandManager",
						"odinRendererStubMessage",
						"ui\\thoughtIcons.xml", -- iconskin
						"explosion", -- icon
						"Bombing Run Requested", -- header text
						"We have requested bombing runs from the Imperial Air Corps", -- text description
						"Right-click to dismiss.", -- action string
						"invasionArrivesResponseStub", -- alert type (for stacking)
						"", -- imagename for bg
						"low", -- importance: low / high / critical
						nil, -- object ID
						60 * 1000, -- duration in ms
						0, state.director) -- snooze
					
				end
				
				send(state.director,"doStage",state.dialogBoxResults)
				return "final"
			elseif state.counter >= 600 then
				
				send("rendCommandManager",
					"odinRendererStubMessage",
					"ui\\thoughtIcons.xml", -- iconskin
					"explosion", -- icon
					"Bombing Run Requested", -- header text
					"Lacking a response, the Imperial Air Corps has decided to initiate bombing runs.", -- text description
					"Right-click to dismiss.", -- action string
					"invasionArrivesResponseStub", -- alert type (for stacking)
					"", -- imagename for bg
					"low", -- importance: low / high / critical
					nil, -- object ID
					60 * 1000, -- duration in ms
					0, state.director) -- snooze
				
				send(state.director,"doStage","stage2_reinforcements")
				return "final"
			end
			
			state.counter = state.counter +1
			return "waiting"
		end,
		
		["suddenly_peace"] = function(state,tags)

			if state.counter == 600 then
			
				send("rendCommandManager",
					"odinRendererFYIMessage",
					state.invaderNationInfo.iconSkin, -- iconskin
					state.invaderNationInfo.iconFlag, -- icon
					"A Diplomatic Solution", -- header text
					"The Invasion is cancelled! A flurry of last-minute negotiations with " .. state.invaderNationInfo.adjective ..
						" diplomats has ended the crisis which prompted invasion. The soldiers of " ..
						state.invaderNationInfo.shortName .. " have been issued a recall order and shall return home without firing a shot. Huzzah!" , -- text description
					"Left-click for more info. Right-click to dismiss.", -- action string
					"invasionCancelled", -- alert type (for stacking)
					"ui\\eventart\\news.png", -- imagename for bg
					"low", -- importance: low / high / critical
					nil, -- object ID
					60 * 1000, -- duration in ms
					0,
					state.director) -- snooze
				
			elseif state.counter <= 0 then
				send("gameSession","setSessionBool","invasionInProgress",false)
				send(state.director,"endEventArc")
				return "final", true
			end
			
			state.counter = state.counter - 1
			
			return "suddenly_peace"
		end,
		
		["final"] = function(state,tags)
			
			return
		end,

		["abort"] = function(state,tags)
			send("gameSession","setSessionBool","invasionInProgress",false)
			send(state.director,"arcAbortDueToRequirements")
			return
		end
	>>
>>
