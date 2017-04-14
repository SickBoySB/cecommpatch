event "planets_align"
<<
	state 
	<<
		int timeout
	     bool alertTriggered
		string dialogBoxResults
	>>

	receive Create( stringstringMapHandle name )
	<< 
		--printl("events","Starvation bailout event started")
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
			printl("events", "Planets Align event firing")
			settimer("Creepy Planets Timer",0)

			send("rendCommandManager",
				"odinRendererTickerMessage",
				"The colony felt a strange dread run through it as the planets aligned!",
				"eldritch",
				"ui\\thoughtIcons.xml")
			
			send("rendCommandManager",
				"odinRendererPlaySoundMessage",
				"alertMad")
			
			state.timeout = 0
			return "real_start"
          end,
          
          ["real_start"] = function(state,tags)               
			incMusic(5,10)
			
			local results = query("gameSpatialDictionary", "allCharactersWithTagRequest", "citizen")
			--Give EVERYONE THE CREEPIES.
			for k,v in pairs(results[1]) do
				local tags = query(v,"getTags")[1]
				if not tags.dead then
					send(v,"hearExclamation", "planetaryAlignment", nil, nil)
				end
			end
			
			send("rendCommandManager",
				"odinRendererChoiceAlertMessage",
				"ui\\thoughtIcons.xml",				-- iconFile
				"eldritch",						-- icon
				"Planetary Alignment",				-- choiceText
				"Click to be Concerned.",			-- choiceTooltipText
				"alert_bar",						-- choiceProgressBarType
				"",								-- choiceDetailsHeader (not used)
				"The colonists shiver as a truly strange celestial event occurs; an alignment of the planets! \z
				A strange blackness spreads out from a single point in the sky, blotting out all other stars. \z
				You feel as though tonight is a very bad night to go outside...",	-- choiceDetailsParagraph
				"ui\\eventart\\monolith_and_stars.png",										-- choiceDetailsArt
				{
					{
						["buttonText"] =  "I'm just going to.. close the windows.",
						["dialogBoxResults"]  = "yes",
						["toolTipText"] = "...and lock them.",
					}
				},						-- choiceTable
				60 * 1000,									-- duration in ms
				SELF,									-- event object (self)
				true,									-- true if event is choice (always true)]]
				nil)									-- director object (null if none)
			
			state.timeout = 0
			return "spawn"
		end,

		["spawn"] = function(state,tags)
		
			local x_max = query("gameSession","getSessionInt","x_max")[1]
			local y_max = query("gameSession","getSessionInt","y_max")[1]

			local spawnLocs = {
				[1] = {x=y_max-32,y=70, w=x_max-32, h=y_max-32},
				[2] = {x=x_max-32,y=70, w=x_max-32, h=y_max-32},
			}
			
			local spawnBox = spawnLocs[ rand(1,#spawnLocs) ]
			local spawnLoc = {x=rand(spawnBox.x, spawnBox.x + spawnBox.w), y=rand(spawnBox.y, spawnBox.y + spawnBox.h)}
			
			local spawnValid = false
			local i = 0
			while not spawnValid do
				local newLoc = gameGridPosition:new()
				newLoc.x = spawnLoc.x
				newLoc.y = spawnLoc.y
				local civ = query( "gameSpatialDictionary", "gridGetCivilization", newLoc )[1]
				local iswater = query( "gameSpatialDictionary", "gridHasSpatialTag", newLoc, "water" )[1]
				printl("events", "planets_align checking if civ or water at pos: " .. tostring(newLoc.x) .. ", " .. tostring(newLoc.y) .. " and civ/water == " .. tostring(civ) .. " / " .. tostring(iswater) )
				if civ < 10*10 then
					-- no! respawn!
					-- make really sure this never tries to spawn off the map.
					-- Theoretically will never happen unless nearly the entire map is covered in buildings, but ...
					spawnBox = spawnLocs[ rand(1,#spawnLocs) ]
					spawnLoc = { x=rand(spawnBox.x -i, spawnBox.x + spawnBox.w +i), y=rand(spawnBox.y -1, spawnBox.y + spawnBox.w +1) }
				elseif iswater then
					-- no! respawn!
					spawnBox = spawnLocs[ rand(1,#spawnLocs) ]
					spawnLoc = { x=rand(spawnBox.x -i, spawnBox.x + spawnBox.w +i), y=rand(spawnBox.y -1, spawnBox.y + spawnBox.w +1) }
				elseif i > 16 then
					-- wow, okay. Just give up.
					printl("events", "WARNING: planets_align seriously couldn't find a place to spawn! Wow.")
					return "final"
				else
					spawnLoc.x = newLoc.x
					spawnLoc.y = newLoc.y
					spawnValid = true
					printl("events", "planets_align : found valid spawn location, placing obeliskians at " .. tostring(spawnLoc.x) .. ", " .. tostring(spawnLoc.y) )
				end
			end
			
               local numEnemies = rand(3,7)
               
               for i=1,numEnemies do
                         
                    local results = query( "scriptManager", "scriptCreateGameObjectRequest",
                                        "obeliskian_group", {
									legacyString = "obeliskian",
									dormant = "false",
									} )
                    
                    local handle = results[1]
                    if handle ~= nil then
                         send(handle, "GameObjectPlace", spawnLoc.x + rand(-4, 4), spawnLoc.y + rand(-4, 4) )
                    end
                    
               end 
			return "final", true
		end,

		["final"] = function(state,tags)  
			return
		end,

		["abort"] = function(state,tags) 
			printl("events", "ominous_dreams: Aborting") 
			return
		end
	>>
>>