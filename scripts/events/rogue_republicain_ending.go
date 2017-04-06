event "rogue_republicain_ending"
<<
	state 
	<<
		bool alertTriggered
		int counter
		string dialogBoxResults
	>>

	receive Create( stringstringMapHandle init )
	<< 
		printl("events","rogue_republicain_ending started")
		state.director = query("gameSession","getSessiongOH","event_director" .. init.director_name)[1]
		state.dialogBoxResults = ""
		
		state.rogue = query(state.director,"getKeyGameobject","rogue")[1]
		state.name = query(state.rogue,"getName")[1]
	>>
	
	receive boxSelectionResult( string id, string message )
	<<
		state.dialogBoxResults = message
	>>

	FSM 
	<<
		["start"] = function(state,tags)
			
			local rogueTags = query(state.rogue,"getTags")[1]
			if rogueTags.dead then
				local icon = "skull"
				local iconskin = "ui\\thoughtIcons.xml"
				
				send("rendCommandManager",
					"odinRendererStubMessage", -- "odinRendererStubMessage",
					iconskin, -- iconskin
					icon, -- icon
					"Rogue Artist Dead", -- header text
					state.name .. " appears to have died. This certainly means that the art won't be happening. On the other hand, perhaps their previous work is more valuable now..?", -- text description
					"Left-click to zoom. Right-click to dismiss.", -- action string
					"republicainRogue", -- alert type (for stacking)
					"ui\\eventart\\travellingPoet.png", -- imagename for bg
					"low", -- importance: low / high / critical
					state.rogue.id, -- object ID
					60 * 1000, -- duration in ms
					0, -- snooze
					state.director)

				return "final"
			end

			send(state.rogue,"makeNeutral")
			send(state.rogue,"addTag","exit_map")
			
			local path_chosen = query(state.director,"getKeyString","path_chosen")[1]
			printl("events","rogue_stahlmarkian_ending path_chosen = " .. tostring(path_chosen))
			
			if path_chosen == "art" then
				
				local result = query(state.director,"getKeyString","result")[1]
				printl("events","rogue_republicain_ending result = " .. tostring(result))
				
				if result == "success" then
                         
                         local successLevel = query(state.director, "getKeyInt", "artPoints")[1]
                         local resultTable = {}
					
                         if successLevel <= 1 then --You are bad. You get nothing.
                             
					     table.insert(resultTable, "none")
						
                         elseif successLevel <= 5 then
						--very minor success
						
                              table.insert(resultTable, "boxed_rare_painting")
                              table.insert(resultTable, "boxed_stained_glass_window")
                              table.insert(resultTable, "boxed_stone_plinth")
                              table.insert(resultTable, "boxed_wood_plinth")
                              table.insert(resultTable, "bricks")
                              table.insert(resultTable, "bricabrac")
                              table.insert(resultTable, "rough_stone_block")
                              table.insert(resultTable, "boxed_wall-mounted_aurochs_head")
						
                         elseif successLevel <= 10 then
						--minor success
						
                              table.insert(resultTable, "boxed_rare_painting")
                              table.insert(resultTable, "boxed_rare_painting")
                              table.insert(resultTable, "boxed_stained_glass_window")
                              table.insert(resultTable, "boxed_stone_plinth")
                              table.insert(resultTable, "boxed_wood_plinth")
                              table.insert(resultTable, "bricks")
                              table.insert(resultTable, "bricabrac")
                              table.insert(resultTable, "rough_stone_block")
                              table.insert(resultTable, "steel_ingots")
                              table.insert(resultTable, "boxed_wall-mounted_aurochs_head")
						
                         elseif successLevel <= 15 then
						-- standard success
						
                              table.insert(resultTable, "boxed_rare_painting")
                              table.insert(resultTable, "boxed_rare_painting")
                              table.insert(resultTable, "boxed_rare_painting")
                              table.insert(resultTable, "boxed_stained_glass_window")
                              table.insert(resultTable, "boxed_stone_plinth")
                              table.insert(resultTable, "boxed_wood_plinth")
                              table.insert(resultTable, "bricks")
                              table.insert(resultTable, "bricabrac")
                              table.insert(resultTable, "rough_stone_block")
                              table.insert(resultTable, "steel_ingots")
                              table.insert(resultTable, "boxed_wall-mounted_aurochs_head")
						
                         elseif successLevel <= 20 then
						-- significant success
						
                              table.insert(resultTable, "boxed_rare_painting")
                              table.insert(resultTable, "boxed_rare_painting")
                              table.insert(resultTable, "boxed_rare_painting")
                              table.insert(resultTable, "boxed_rare_painting")
                              table.insert(resultTable, "boxed_rare_painting")
                              table.insert(resultTable, "boxed_stained_glass_window")
                              table.insert(resultTable, "boxed_stone_plinth")
                              table.insert(resultTable, "boxed_wood_plinth")
                              table.insert(resultTable, "bricks")
                              table.insert(resultTable, "steel_ingots")
                              table.insert(resultTable, "boxed_wall-mounted_aurochs_head")
						
                         else
						--super success!!
                              table.insert(resultTable, "boxed_rare_painting")
                              table.insert(resultTable, "boxed_rare_painting")
                              table.insert(resultTable, "boxed_rare_painting")
                              table.insert(resultTable, "boxed_rare_painting")
                              table.insert(resultTable, "boxed_rare_painting")
                              table.insert(resultTable, "boxed_rare_painting")
                              table.insert(resultTable, "boxed_rare_painting")
                              table.insert(resultTable, "boxed_rare_painting")
                              table.insert(resultTable, "boxed_grenade_launcher_locker") --!??!?!?
                              table.insert(resultTable, "boxed_stone_plinth")
                              table.insert(resultTable, "boxed_wood_plinth")
                              table.insert(resultTable, "steel_ingots")
                              table.insert(resultTable, "boxed_wall-mounted_aurochs_head")
                         end

                         state.itemPick = resultTable[rand(1,#resultTable)]
					
                         local mysteriousResult = "The work of art appears to be something truly mysterious. You're not quite sure what it is, honestly. You blink, and suddenly it's gone. Huh. You feel like you've been had."
                         
                         if state.itemPick == "boxed_rare_painting" then
						
                              mysteriousResult = "The work of art is an astonishing painting, capturing your colony in glorious, swirling colors - truly a masterwork! \z
                              Your colonists' hearts swell with pride."
                              --todo: give positive memory to all colonists
						
						if not query("gameSession","getSessionBool","patronOfTheArts")[1] then
							send("gameSession","setSessionBool","patronOfTheArts",true)
							send("gameSession", "setSteamAchievement", "patronOfTheArts")
						end
						
                         elseif state.itemPick == "boxed_grenade_launcher_locker" then
						
                              mysteriousResult = "The artist, mysteriously, presents you with an incredibly dangerous looking weapon. 'Thank you for allowing me to complete my project in peace. \z
                              I'm sorry for deceiving you - the Republique desires my unique skills, and it would be best if I wasn't found. Use it as you will - I have learned all I can from it.' \z
                              You've no idea what to think of this turn of events."
                              --todo: Huge hit to republique relations.
						
                         elseif state.itemPick == "boxed_stained_glass_window" then
						
                              mysteriousResult = "The work of art is beautiful, although not quite what you imagined; a stained glass window! It would look rather nice in a chapel, you think."
						
                         elseif state.itemPick == "boxed_stone_plinth" or state.itemPick == "boxed_wood_plinth" then
						
                              mysteriousResult = "The 'work of art' appears to be a rather amateurish plinth. You suspect you've just been had by an art student. Curses!"
						
                         elseif state.itemPick == "steel_ingots" then
						
                              mysteriousResult = "'My plan has failed!' declares the artist. 'I cannot work such ridiculous material. I am ruined!' \z
                              They leave behind a pile of steel ingots - were they really trying to \z
                              make art out of that...?"
						
                         elseif state.itemPick == "boxed_wall-mounted_aurochs_head" then
						
                              mysteriousResult = "The 'work of art' appears to be a mounted aurochs head. 'Hunted it myself!' the rogue proclaims. \z
                              You suppose taxidermy is a form of art, but you can't help but be a bit disappointed."
						
                         elseif state.itemPick == "bricks" or state.itemPick == "rough_stone_block" then
						
                              mysteriousResult = "The 'work of art' appears to be a block of stone. Is this... some new 'art' thing? \z
                              You're afraid to touch it in case it's incredibly valuable. A labourer comes by and hauls it to the stockpile."
						
                         elseif state.itemPick == "bricabrac" then
						
                              mysteriousResult = "The 'work of art' appears to be a completely random pile of nonsense. As far as you're concerned, this definitely isn't art. \z
                              It might be useful for making some decor of your own, you suppose."
						
                         end
                         
					local icon = "decor_uc_icon"
					local iconskin = "ui\\orderIcons.xml"
                         
                         --Let's assemble the art string!
                         local artObjectCount = query(state.director, "getKeyInt", "artObjectCount")[1]
                         local randomArtString = ""
                         local randArtNameTable = {
                              "Piles",
                              "Sadness",
                              "Glop",
                              "Mystery",
                              "Goading",
                              "Beholdment",
                              "Frippery",
                              "Constabulary",
                              "Empirical",
                              "Goosed",
                              "Withering",
                              "Failure",
                              "Lordship",
                              "Cruelty",
                              "Aurochs",
                              "Deer",
                              "Beetles",
                              "Turtles",
                              }
                         local randArtPrefixTable = {
                              "The Dregs of ",
                              "The Lords of ",
                              "The Sudden Appearance of",
                              "Whence ",
                              "And then ",
                              "Forsooth, ",
                              "I wish I had a",
                              "The ",
                              "The Great ",
                              "Then, I beheld ",
                              "My Favorite ",
                              "A Study on ",
                              "The Murder of ",
                              "The Beating of ",
                              "The Cruelty of ",
                              "The Kindness of ",
                              "The Beauty of ",
                              "The Unending Horror of ",
                              "A Portrait of ",
                              "Suddently, ",
                              "In the End, All Is ",
                              }
                         local randArtMidfixTable = {
                              " and ",
                              " and ",
                              " and ",
                              " with ",
                              " with ",
                              " near ",
                              " near ",
                              ", yes, ",
                              ", yes, ",
                              " but also ",
                              " containing ",
                              " withheld from ",
                              " nearby ",
                              " juxtaposed against ",
                              " in contemplation of ",
                              "! And "
                              }
                         local randArtSuffixTable = {
                              " for ever",
                              ", to my astonishment",
                              "!",
                              ", alas!",
                              ", I do say",
                              ", perhaps",
                              " as such",
                              " in my dreams",
                              " by the sea",
                              " as metaphor",
                              " en plein air",
                              ", disgusting!",
                              ", oh..."
                              }
                         
                         --[[A VERY BRIEF POLICY DESCRIPTION FOR THIS MESS:
                         * When we add a word we either add a stock word or a "nearby object" word.
                         * If it's a nearby object word, we have a 1/2 chance of shortening it to either the first or last word in the string.
                         If it's shortened, the next word in line will also be shortened to allow them to mash together nicely.
                         If it's not shortened, we add a midfix (" and, "  ", with" to make it pair nicely)
                         
                         
                         ]]
                         local setMidFix = false
                         local shortened = false
                         if rand(1,4) == 1 then
                              randomArtString = randomArtString .. randArtPrefixTable[rand(1,#randArtPrefixTable)]
                         end
                         if artObjectCount > 0 and rand(1,3) > 1 then
                              local artObject = query(state.director, "getKeyString", "artObject1")[1]
                              if rand(1,2) == 1 then
                                   if rand(1,2) == 1 then --pick the first or last word randomly
                                        artObject = string.match(artObject,"(%a+)")
                                   else
                                        local sub = string.gsub(artObject,"(%a+)",
                                             function(w)
                                             artObject = w
                                             return nil
                                             end)
                                        --brute force sets artObject to last word in string bc I am a bad programmer
                                   end
                                   shortened = true
                              end
                              randomArtString = randomArtString .. artObject
                              setMidFix = false
                              if shortened == false then
                                   randomArtString = randomArtString .. randArtMidfixTable[rand(1,#randArtMidfixTable)]
                                   setMidFix = true
                              end
                         else
                              randomArtString = randomArtString .. randArtNameTable[rand(1,#randArtNameTable)]
                              setMidFix = false
                              if rand(1,3) > 1 then
                                   randomArtString = randomArtString .. randArtMidfixTable[rand(1,#randArtMidfixTable)]
                                   setMidFix = true
                              else
                                   shortened = true
                              end
                         end
                         
                         if artObjectCount > 1 and rand(1,4) > 1 then
                              local artObject = query(state.director, "getKeyString", "artObject2")[1]
                              if (shortened == true) or (rand(1,2) == 1) then
                                   if rand(1,2) == 1 then --pick the first or last word randomly
                                        artObject = string.match(artObject,"(%a+)")
                                   else
                                        local sub = string.gsub(artObject,"(%a+)",
                                             function(w)
                                             artObject = w
                                             return nil
                                             end)
                                        --brute force sets artObject to last word in string bc I am a bad programmer
                                   end
                                   if shortened == true then --This silly bit of code makes sure we properly flip the case if this was a bookend
                                        artObject = string.lower(artObject) --If this is the 2nd half of a mash, lowercase it. 
                                        shortened = false
                                   elseif shortened == false then
                                        shortened = true
                                   end
                              end
                              randomArtString = randomArtString .. artObject
                              setMidFix = false
                              if rand(1,3) > 1 then
                                   randomArtString = randomArtString .. randArtMidfixTable[rand(1,#randArtMidfixTable)]
                                   setMidFix = true
                              end
                         elseif rand(1,4) == 1 then
                              randomArtString = randomArtString .. randArtNameTable[rand(1,#randArtNameTable)]
                              setMidFix = false
                              if rand(1,3) > 1 then
                                   randomArtString = randomArtString .. randArtMidfixTable[rand(1,#randArtMidfixTable)]
                                   setMidFix = true
                              else
                                   shortened = true
                              end
                         end
                         
                         if artObjectCount > 2 and ((rand(1,3) > 1) or (setMidFix == true)) then
                              local artObject = query(state.director, "getKeyString", "artObject3")[1]
                              if (shortened == true) or (rand(1,2) == 1) then
                                   if rand(1,2) == 1 then --pick the first or last word randomly
                                        artObject = string.match(artObject,"(%a+)")
                                   else
                                        local sub = string.gsub(artObject,"(%a+)",
                                             function(w)
                                             artObject = w
                                             return nil
                                             end)
                                        --brute force sets artObject to last word in string bc I am a bad programmer
                                   end
                                   if shortened == true then
                                        artObject = string.lower(artObject) --If this is the 2nd half of a mash, lowercase it. 
                                   end
                              end
                              randomArtString = randomArtString .. artObject
                              setMidFix = false
                         end
                         
                         if setMidFix == true then
                              randomArtString = randomArtString .. randArtNameTable[rand(1,#randArtNameTable)]
                         end
                         
                         if rand(1,3) == 1 then
                              randomArtString = randomArtString .. randArtSuffixTable[rand(1,#randArtSuffixTable)]
                         end
					
                         
                         if state.itemPick ~= "none" then --spawn an item
                             
					    local handle = query("scriptManager",
									 "scriptCreateGameObjectRequest",
									 "item",
									 {legacyString = state.itemPick,
                                               displayNameOverride = randomArtString})[1]--*ENDLESS LAUGHING*
						
                              local pos = query(state.rogue, "gridGetPosition")[1]
                              send(handle,"ClaimItem")
                              send(handle, "GameObjectPlace", pos.x, pos.y )
                              send("rendCommandManager",
                              "odinRendererPlaySoundMessage",
                              "alertGood")
                              
                              send("rendCommandManager",
                                   "odinRendererFYIMessage",
                                   iconskin, -- iconskin
                                   icon, -- icon
                                   "Art Complete!", -- header text
                                   state.name .. " has finished their work, which they have entitled '" .. randomArtString .. "'. " .. mysteriousResult, -- state.gangName .. " has broken our truce, and intends to resume hostilities!",
                                   "Left-click to learn more, Right-click to dismiss.", -- tooltip string
                                   "stahlmarkianRogueEnd", -- alert type (for stacking)
                                   "ui//eventart//prosperity_wall_poster.png", -- imagename for bg
                                   "low", -- importance: low / high / critical
                                   handle, -- state.renderHandle, -- object ID
                                   3 * 60 * 1000, -- duration in ms
                                   0, -- "snooze" time if triggered multiple times in rapid succession
                                   nil)
                         else
                              send("rendCommandManager",
                                   "odinRendererFYIMessage",
                                   iconskin, -- iconskin
                                   icon, -- icon
                                   "Art Complete!", -- header text
                                   state.name .. " has finished their work, which they have entitled '" .. randomArtString .. "'. " .. mysteriousResult, -- state.gangName .. " has broken our truce, and intends to resume hostilities!",
                                   "Left-click to learn more, Right-click to dismiss.", -- tooltip string
                                   "stahlmarkianRogueEnd", -- alert type (for stacking)
                                   "ui//eventart//prosperity_wall_poster.png", -- imagename for bg
                                   "low", -- importance: low / high / critical
                                   state.rogue.id, -- state.renderHandle, -- object ID
                                   3 * 60 * 1000, -- duration in ms
                                   0, -- "snooze" time if triggered multiple times in rapid succession
                                   nil)
                         end
                         
                         
					
				elseif result == "combat_abort" then
					
                         send("rendCommandManager",
                              "odinRendererPlaySoundMessage",
                              "alertNeutral")
                         
					local icon = "retreat"
					local iconskin = "ui\\thoughtIcons.xml"
					
					send("rendCommandManager",
						"odinRendererStubMessage", -- "odinRendererStubMessage",
						iconskin, -- iconskin
						icon, -- icon
						"Republicain Artist Flees", -- header text
						state.name .. " is fleeing from this dangerous colony! It is safe to assume that we won't be receiving any artwork.", -- text description
						"Left-click to zoom. Right-click to dismiss.", -- action string
						"republicainRogueEnd", -- alert type (for stacking)
						"ui\\eventart\\travellingPoet.png", -- imagename for bg
						"low", -- importance: low / high / critical
						state.rogue.id, -- object ID
						60 * 1000, -- duration in ms
						0, -- snooze
						state.director)
					
					send(state.rogue,"addTag","exit_map_run")
				end
				
			elseif path_chosen == "nope" then
				-- just leave.
				local icon = "decor_uc_icon"
				local iconskin = "ui\\orderIcons.xml"
				
				send("rendCommandManager",
					"odinRendererStubMessage", -- "odinRendererStubMessage",
					iconskin, -- iconskin
					icon, -- icon
					"Republicain Artist Leaving", -- header text
					state.name .. " is leaving our colony, the offer made having been denied.", -- text description
					"Left-click to zoom. Right-click to dismiss.", -- action string
					"republicainRogueEnd", -- alert type (for stacking)
					"ui\\eventart\\travellingPoet.png", -- imagename for bg
					"low", -- importance: low / high / critical
					state.rogue.id, -- object ID
					60 * 1000, -- duration in ms
					0, -- snooze
					state.director)
                    
                    send("rendCommandManager",
                              "odinRendererPlaySoundMessage",
                              "alertNeutral")
				
			end
			return "final"
		end,
		
		["final"] = function(state,tags)
			send(state.director,"endEventArc")
			return
		end,

		["abort"] = function(state,tags)
			send(state.director,"endEventArc")
			return
		end
	>>
>>
