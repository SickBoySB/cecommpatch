event "empire_times_ending"
<<
	state 
	<<
		bool alertTriggered
		int counter
		string dialogBoxResults
	>>

	receive Create( stringstringMapHandle init )
	<< 
		printl("events","empire_times_ending started")
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
					"Article Writer Dead", -- header text
					state.name .. " appears to have died - this certainly means there will be no forthcoming article from the Empire Times. Hopefully no one traces this back to your colony.", -- text description
					"Left-click to zoom. Right-click to dismiss.", -- action string
					"republicainRogue", -- alert type (for stacking)
					"ui\\eventart\\news.png", -- imagename for bg
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
			
			if (path_chosen == "yes") or (path_chosen == "nope") then
				
				local result = query(state.director,"getKeyString","result")[1]
				printl("events","rogue_republicain_ending result = " .. tostring(result))
				
				if result == "success" then
                         
					--HERE WE GO
					local columnObjectCount = query(state.director, "getKeyInt", "columnObjectCount")[1]
					local columnObjects = {}
					if columnObjectCount > 0 then
						for i=1,columnObjectCount do
							table.insert(columnObjects, query(state.director, "getKeyString", "columnObject"..i)[1])
						end
					end
					local columnBuildingCount = query(state.director, "getKeyInt", "columnBuildingCount")[1]
					local columnBuildings = {}
					local columnBuildingQualities = {}
					if columnBuildingCount > 0 then
						for i=1,columnBuildingCount do
							table.insert(columnBuildings, query(state.director, "getKeyString", "articleBuilding"..i)[1])
							table.insert(columnBuildingQualities, query(state.director, "getKeyInt", "articleBuildingQuality"..i)[1])
						end
					end
					local columnPersonCount = query(state.director, "getKeyInt", "columnPersonCount")[1]
					local columnPersons = {}
					local columnPersonProfs = {}
					local columnPersonMoods = {}
					if columnPersonCount > 0 then
						for i=1,columnPersonCount do
							table.insert(columnPersons, query(state.director, "getKeyString", "columnPerson"..i)[1])
							table.insert(columnPersonProfs, query(state.director, "getKeyString", "columnPersonProf"..i)[1])
							table.insert(columnPersonMoods, query(state.director, "getKeyString", "columnPersonMood"..i)[1])
						end
					end
                         
					local icon = "decor_uc_icon"
					local iconskin = "ui\\orderIcons.xml"
                         
                         --Alright, how are we doing this.
                         local article = ""
					local positivity = 5 --we'll.. use this to measure how nice the article is, I guess.
					local columnType = query(state.director,"getKeyString","column_type")[1]
					local reporterType = query(state.director,"getKeyString","reporter_type")[1]
					--printl("CHRIS", "Reportertype is " .. reporterType .. " and column type is " .. columnType)
					
					if path_chosen == "nope" then 
						if columnType ~= "murder mystery" then
							positivity = positivity - 2 --Rebellious reporters have little use for you!
						else
							positivity = positivity + 2 --...unless they're murder mystery writers in which case they relish the thrill.
						end
					end
					
					local positivityTable = { --This defines the reporter's reaction to things based on their mood
						"horrifying", --1
						"depressing", --2
						"quite unfortunate", --3
						"lacking", --4
						"tolerable", --5
						"acceptable", --6
						"lovely", --7
						"excellent", --8
						"wonderful", --9
						"top-shelf" --10
					}
					
					local positivityDescTable = { --Used to describe the result at the end.
						"very negative", --1
						"negative", --2
						"negative", --3
						"negligible", --4
						"negligible", --5
						"negligible", --6
						"positive", --7
						"positive", --8
						"very positive", --9
						"very positive" --10
					}
					
					
					local interestTable = {} --these strings are things they are looking for.
					interestTable["cultural interest"] = {
						"fascinating local events",
						"the local cuisine",
						"colonial traditions",
						"this colony's unique crafts"
					}
					interestTable["'fascinating people of the Colonies'"] = {
						"whoever happened to be about",
						"the most interesting person here",
						"the colony manager",
					}
					interestTable["cooking"] = {
						"the typical meal",
						"the local cuisine",
						"secret family recipes",
					}
					interestTable["gardening"] = {
						"everyone's favorite crop, Cabbages",
						"the local farms",
						"colonial gardening practises",
					}
					interestTable["arts and crafts"] = {
						"colonial decor practises",
						"fascinating local crafts", 
						"the local art scene",
					}
					interestTable["dramatic adventure"] = {
						"the unseen wonders of the Frontier",
						"the deadly animals of this Foreign Land", 
						"the secrets of this hostile waste",
					}
					interestTable["murder mystery"] = {
						"the answer to last week's Mystery",
						"just who, in fact, dunit", 
						"the blood-soaked truth of this violent Den of Ruffians",
					}
					
					local simpleInterestTable = {} --these strings are things they are looking for.
					simpleInterestTable["cultural interest"] = {
						"culture",
						"lifestyle",
						"traditions",
						"crafts"
					}
					simpleInterestTable["'fascinating people of the Colonies'"] = {
						"people",
						"luminaries",
						"colony manager",
					}
					simpleInterestTable["cooking"] = {
						"cooking",
						"cuisine",
						"recipes",
					}
					simpleInterestTable["gardening"] = {
						"crops",
						"farms",
						"gardening practises",
					}
					simpleInterestTable["arts and crafts"] = {
						"decor",
						"crafts", 
						"art",
					}
					interestTable["dramatic adventure"] = {
						"unseen secrets",
						"deadly animals", 
						"harsh struggle to survive",
					}
					simpleInterestTable["murder mystery"] = {
						"mystery",
						"clues", 
						"truth",
					}
					
					local emotionTable = {} -- Concepts the writer enjoys, used to describe things positively (usually)
					emotionTable["fast-talking"] = { --sort of a con man
						"savvy",
						"quick-thinking",
						"double-dealing",
						"well-read",
						"speedy"
					}
					emotionTable["quiet, purposeful"] = { --stoic, refined
						"dour",
						"tolerant",
						"stoic",
						"intelligent",
						"patient"
					}
					emotionTable["oddly meticulous"] = { --a robot???
						"organised",
						"contemporous",
						"thorough",
						"clear-eyed",
						"orderly"
					}
					emotionTable["scowling"] = { --hates everyone including their readers
						"dismal",
						"depressing",
						"unfortunate",
						"useless",
						"intolerable"
					}
					emotionTable["frantic"] = { --eager to please
						"likely",
						"lovely",
						"alert",
						"great",
						"eager"
					}
					emotionTable["grinning"] = { --uses overly familiar language
						"silly",
						"dear",
						"friendly",
						"most excellent",
						"kind"
					}
					emotionTable["pleasant"] = { --very drab
						"nice",
						"quite alright",
						"pleasant",
						"gentle",
						"proper"
					}
					
					local columnNames = {} --Name of the column.
					columnNames["cultural interest"] = {
						"Fascinating Locales",
						"Unusual Eats" ,
						"Best Traditions",
						"Craftswork of the Empire"
					}
					columnNames["'fascinating people of the Colonies'"] = {
						"Peoples All About The World",
						"Bright Minds of the Empire",
						"Our Finest Leaders",
					}
					columnNames["cooking"] = {
						"Simple Cooking",
						"Fine Empire Cuisine",
						"The Secrets of Taste",
					}
					columnNames["gardening"] = {
						"Cabbages! Why Do We All Love Them So Much?",
						"Best Farming Practises",
						"Your Garden, My God!",
					}
					columnNames["arts and crafts"] = {
						"Modest Interior Decoration",
						"What Might This Thing Be?", 
						"But Is It Art?",
					}
					columnNames["dramatic adventure"] = {
						"267 Astonishing Places You Aren't Visting",
						"Horrible Creatures and Their Behaviours", 
						"The Worst Places To Go",
					}
					columnNames["murder mystery"] = {
						"Weekly Murders!, Weekly Edition",
						"Tracking The Culprit", 
						"The Gory Truth (Not For Those Under The Age of 25)",
					}
					
					local Starts = {
						"Welcome back, " .. emotionTable[reporterType][rand(1,#emotionTable[reporterType])] .. " readers, \z
						to another installment of " .. columnNames[columnType][rand(1,#columnNames[columnType])] .. " with " .. state.name .. ".",
						
						"Hello " .. emotionTable[reporterType][rand(1,#emotionTable[reporterType])] .. " readers. I, " .. state.name .. ", have the pleasure of presenting another installment of '" .. columnNames[columnType][rand(1,#columnNames[columnType])] .. "'.",
						
						"In this week's column, we return to: '" .. columnNames[columnType][rand(1,#columnNames[columnType])] .. "' with " .. state.name .. ".",
						
						emotionTable[reporterType][rand(1,#emotionTable[reporterType])]:gsub("^%l", string.upper) .. " readers, this is '" .. columnNames[columnType][rand(1,#columnNames[columnType])] .. "'. I, as always, am " .. state.name .. "."
					}
					
					--Where are we?
					local biome = query("gameSession", "getSessionString", "biome")[1]
					--The introduction should also mention why they're at a colony. Setup their basic inclination toward positivity based on subject too.
					local Introduction = ""
					if (columnType == "dramatic adventure") or (columnType == "murder mystery") then --They like exciting biomes.
						local horribleThing = ""
						if (columnType == "dramatic adventure") then
							horribleThing = "Beast's fangs"
						else
							horribleThing = "Murderer's dagger"
						end
						
						local locationIndicator = ""
						local locationAdjective = ""
						if biome == "cold" then
							locationIndicator = "in the heart of winter, where the cold bites deeper than any " .. horribleThing
							locationAdjective = "desolate land"
							positivity = positivity + 1
						elseif biome == "desert" then
							locationIndicator = "under the endless sands, where the exhaustion and quivering air can hide a waiting " .. horribleThing
							locationAdjective = "cruel waste"
							positivity = positivity + 1
						elseif biome == "tropical" then
							locationIndicator = "deep in the harsh Jungles, where a deadly plant's poison can kill you as quickly as a " .. horribleThing
							locationAdjective = "forgotten valley"
							positivity = positivity + 1
						else --temperate
							locationIndicator = "in the least likely place imaginable; an idyllic colony in the calm flatlands, far from the reach of any " .. horribleThing
							locationAdjective = ""
							positivity = positivity - 1
						end
						Introduction = "We begin our story " .. locationIndicator .. "; Upon arriving in this " .. locationAdjective ..", I "
					elseif (columnType == "cultural interest") or (columnType == "'fascinating people of the Colonies'") then
						positivity = positivity + 1 --They're a colony interest column, so they're generally happy investigating them
						local locationIndicator = ""
						if biome == "cold" then
							locationIndicator = "rather frigid and unpleasant unpleasant place; a colony deep in the colds of the north"
						elseif biome == "desert" then
							locationIndicator = "colony in the deserts of the New Continent, where survival is a harsh daily excercise"
						elseif biome == "tropical" then
							locationIndicator = "settlement deep in the jungles; these brave colonists toil daily to extract valuable crops to export back to the Colonies"
							positivity = positivity + 1
						else --temperate
							locationIndicator = "settlement in the temperate hills; truly a wonderful chance to sample the bounties of the New Continent's breadbasket"
							positivity = positivity + 1
						end
						Introduction = "Continuing our ongoing investigation of the Colonies, we come now to a " .. locationIndicator.. ". Upon my arrival, I "
					else -- A non-colonial interest column.
						local locationIndicator = ""
						if biome == "cold" then
							positivity = positivity - 1
							locationIndicator = "in a rather cold, horrid place: Alas, I'm told it was the only location available"
						elseif biome == "desert" then
							positivity = positivity - 1
							locationIndicator = "in a terrible desert - why anyone would settle here is a mystery to me. The heat strained me daily, but I perservered for you, my " .. emotionTable[reporterType][rand(1,#emotionTable[reporterType])] .. " readers"
						elseif biome == "tropical" then
							locationIndicator = "in a disgusting, wet jungle; how anyone tolerates this place I do not know"
							positivity = positivity - 1
						else --temperate
							locationIndicator = "in a surprisingly tolerable hilly area; while the accomodations were lacking, the weather at least was quite fine"
							positivity = positivity + 1
						end
						Introduction = "While this column generally covers Empire interest topics, this week I had the rare opportunity to cover one of Her Majesty's Colonies; a settlement ".. locationIndicator ..". Upon arriving at the colony, I "
					end
					
					local randintroTable = {
						"immediately set to work finding a source of information on ",
						"began to look into the subject at hand, ",
						"searched out a subject to begin our study of ",
						"could sense this place would be bursting with gossip on ",
						"began a subtle and careful investigation into "
						}
					
					Introduction = Introduction .. randintroTable[rand(1,#randintroTable)] .. interestTable[columnType][rand(1,#interestTable[columnType])] .. ". "
					
					if (reporterType == "frantic") or (reporterType == "pleasant") then
						positivity = positivity + 1
					elseif (reporterType == "scowling") then
						positivity = positivity - 1
					end
					local numSubjects = 0
					local break1 = ""
					local break2 = ""
					local break3 = ""
					local Subject1 = ""
					--Okay, now the article is introduced. Next, the subject matter.
					if (columnPersonCount > 0) and (rand(1,4) > 1) then
						numSubjects = numSubjects + 1
						break1 = "\n\n"
						local pickPerson = rand (1,columnPersonCount)
						
						local randMiddleTable = {
							"I soon found ",
							"I quickly sighted ",
							"Upon arrival I had spotted my old acquaintance ",
							"After applying a bit of the old charm, I was referred to ",
							"After a careful canvas of the area, I chose to approach "
						}
						local randLateMidTable = {
							", who on first approach seemed a " .. emotionTable[reporterType][rand(1,#emotionTable[reporterType])] .. " sort. ",
							", who seemed an altogether unusual " .. columnPersonProfs[pickPerson] .. ". ",
							", a rather " .. emotionTable[reporterType][rand(1,#emotionTable[reporterType])] .. " personage, I must say. ",
							", a rather " .. emotionTable[reporterType][rand(1,#emotionTable[reporterType])] .. " " .. columnPersonProfs[pickPerson] .. ", I must say. ",
							", and soon learned they were quite the " .. emotionTable[reporterType][rand(1,#emotionTable[reporterType])] .. " type. ",
							", and soon learned they were a " .. columnPersonProfs[pickPerson] .. " - fascinating! ",
							", against my better judgement - I typically avoid the " .. emotionTable[reporterType][rand(1,#emotionTable[reporterType])] .. " sort. ",
							", against my better judgement - I typically avoid those of the " .. columnPersonProfs[pickPerson] .. " profession. "
						}
						local randReplyTable = {
							"We had quite a discussion, and I must say - the things I heard ",
							"As for what they had to say, their words ",
							"I soon learned a great many facts, some of which quite "
						}
						local randAdjectiveTable ={
							"astonished me: ",
							"shocked me: ",
							"confirmed my suspicions: ",
							"were quite as expected: ",
							"were rather droll: ",
							"confounded every expectation: ",
							"fascinated me: "
						}
						local randColonyResultTable ={
							"it appears this colony is quite a ",
							"the truth of this colony is that of a ",
							"I learned that this colony is a particularly ",
							"I was told tales painting this colony as a most ",
						}
						local reactionTable = {}
						if columnPersonProfs[pickPerson] == "fear" then
							table.insert(reactionTable, "terrifying")
							table.insert(reactionTable, "horrific")
							table.insert(reactionTable, "concerning")
							table.insert(reactionTable, "worrying")
							table.insert(reactionTable, "terrifying")
							positivity = positivity - 1
						elseif columnPersonProfs[pickPerson] == "anger" then
							table.insert(reactionTable, "grating")
							table.insert(reactionTable, "irritating")
							table.insert(reactionTable, "frustrating")
							table.insert(reactionTable, "fruitless")
							table.insert(reactionTable, "angering")
							positivity = positivity - 1
						elseif columnPersonProfs[pickPerson] == "despair" then
							table.insert(reactionTable, "lonely")
							table.insert(reactionTable, "upsetting")
							table.insert(reactionTable, "twisted")
							table.insert(reactionTable, "cruel")
							table.insert(reactionTable, "miserable")
							positivity = positivity - 1
						else --happy
							table.insert(reactionTable, "lovely")
							table.insert(reactionTable, "enjoyable")
							table.insert(reactionTable, "satisfying")
							table.insert(reactionTable, "celebration-worthy")
							table.insert(reactionTable, "joyful")
							positivity = positivity + 1
						end
						local randReactionReactionTable ={
							reactionTable[rand(1,#reactionTable)] .. " place! The " .. reactionTable[rand(1,#reactionTable)] .. " tales I was told were " .. reactionTable[rand(1,#reactionTable)] .. ", indeed.",
							reactionTable[rand(1,#reactionTable)] .. " sort of place to Settle. I found myself expecting to see an event most " .. reactionTable[rand(1,#reactionTable)] .. " at any moment.",
							reactionTable[rand(1,#reactionTable)] .. " land. Truly, I never expected to find a place described in such " .. reactionTable[rand(1,#reactionTable)] .. " terms."
							--I'll.. think of more later?
						}
						
						Subject1 = randMiddleTable[rand(1,#randMiddleTable)] .. columnPersons[pickPerson] .. randLateMidTable[rand(1,#randLateMidTable)] ..
						randReplyTable[rand(1,#randReplyTable)] .. randAdjectiveTable[rand(1,#randAdjectiveTable)] .. randColonyResultTable[rand(1,#randColonyResultTable)] .. randReactionReactionTable[rand(1,#randReactionReactionTable)]  .. " "
						--Format is: I found  .. person .. who seemed (adjective). .. They told me .. something surprising .. this colony is .. Thing! .. indeed.
					end
					local section1to2Bridge = ""
					if Subject1 ~= "" then
						section1to2Bridge = "With that information in-mind, I "
					end
					
					local Subject2 = ""
					if (columnBuildingCount > 0) then
						numSubjects = numSubjects + 1
						break2 = "\n\n"
						local pickBuilding = rand (1,columnBuildingCount)
						local randMiddleTable = {
							"soon found my way to the ",
							"quickly spotted the ",
							"saw a likely place across the way, known as ",
							"put a bit of pressure to the local populace, and was referred to the nearby ",
							"consulted my research; it suggested the ideal location for a scoop would be a place known as the "
						}
						if Subject1 ~= "" then
							table.insert(randMiddleTable, "took the information I had gleaned from my previous interrogation and headed to the ")
							table.insert(randMiddleTable, "knew I must be hot on the trail. I made my next stop the ")
							table.insert(randMiddleTable, "had the lay of the land firmly set in my mind. I knew my next stop must be none other than the ")
							table.insert(randMiddleTable, "was well situated in the colony's social groups. My new friends whispered that my next stop should be the ")
						end
						
						local randMiddleTable2 = {
							", a rather " .. positivityTable[positivity] .. "-looking place, I must say. ",
							", a " .. positivityTable[positivity] .. " little building for a place such as this. ",
							", which seemed quite " .. positivityTable[positivity] .. " at first glance. ",
							". My first impression was of a rather " .. positivityTable[positivity] .. " place. ",
						}
						
						local randDeterminationTable = {
							"I immediately carried out an investigation of the premises and determined it was ",
							"Upon further investigation, I slowly began to realize this location was ",
							"Against my better judgement, I looked into the locale and soon learned that it was "
						}
						local buildingAssessment = ""
						if columnBuildingQualities[pickBuilding] > 4 then --Super good
							buildingAssessment = "rather wonderful"
							positivity = positivity + 1
						elseif columnBuildingQualities[pickBuilding] > 0 then --good
							buildingAssessment = "a rather nice place"
							positivity = positivity + 1
						elseif columnBuildingQualities[pickBuilding] == 0 then --Meh
							buildingAssessment = "quite unremarkable"
						elseif columnBuildingQualities[pickBuilding] < -4 then --Super bad
							buildingAssessment = "dismal and horrid"
							positivity = positivity - 1
						elseif columnBuildingQualities[pickBuilding] < 0 then --bad
							buildingAssessment = "rather unfortunate"
							positivity = positivity - 1
						end
						local randSentenceEndTable = {
							", indeed. What this implies about the Colony's bureaucrat is best left up to the reader's imagination.",
							", I must say. One must only assume this reflects on the Colony at large.",
							" - a clear Indicator of the general state of this Colony."
						}
						Subject2 = randMiddleTable[rand(1,#randMiddleTable)] .. columnBuildings[pickBuilding] ..
						randMiddleTable2[rand(1,#randMiddleTable2)] .. randDeterminationTable[rand(1,#randDeterminationTable)] ..
						buildingAssessment .. randSentenceEndTable[rand(1,#randSentenceEndTable)]
					end
					
					local section2to3Bridge = ""
					if numSubjects == 0 then
						section2to3Bridge = "Lacking any other leads, I decided to investigate the nearby area. I "
					elseif numSubjects == 1 then
						section2to3Bridge = "With this information, I knew exactly what would yield the final bit of info I needed. I "
					end
					local Subject3 = ""
					if (columnObjectCount > 0) and (numSubjects < 2) then
						numSubjects = numSubjects + 1
						break3 = "\n\n"
						local pickObject = rand (1,columnObjectCount)
						positivity = positivity + 1
						local randMiddleTable = {
							"soon found a ",
							"quickly spotted a ",
							"saw a likely object nearby, a ",
							"determined empircally that the most likely source of information would be the nearby ",
							"listened in surreptitously and discovered the existince of a nearby "
						}
						local randMiddleTable2 = {
							". This object ",
							". I couldn't help but feel this ",
							". I was pleasantly surprised to find it ",
							". After careful examination, I determined it ",
							". Surprised, I found that it "
						}
						local transitionDesc = {}
						if columnType == "murder mystery" then
							table.insert(transitionDesc, "pointed me directly to ")
							table.insert(transitionDesc, "revealed to me the secret of ")
							table.insert(transitionDesc, "yielded all the clues I needed to learn ")
						elseif columnType == "dramatic adventure" then
							table.insert(transitionDesc, "was truly a deadly device; a clear example of ")
							table.insert(transitionDesc, "was full of mystery; one of ")
							table.insert(transitionDesc, "most certainly contained deep secrets; objects like this are the very reason I sought out ")
						else
							table.insert(transitionDesc, "was an excellent example of ")
							table.insert(transitionDesc, "was a clear indicator of ")
							table.insert(transitionDesc, "unquestionably was a fine example of ")
						end
						Subject3 = randMiddleTable[rand(1,#randMiddleTable)] .. columnObjects[pickObject] .. randMiddleTable2[rand(1,#randMiddleTable2)] .. transitionDesc[rand(1,#transitionDesc)] .. columnNames[columnType][rand(1,#columnNames[columnType])] .. "."
					elseif numSubjects < 2 then --Make something up!
						positivity = positivity - 1
						local randExcuseTable = {
							"I was unable to find what I sought. It seems this colony would defy my quest for information to the very end. ",
							"I came down with illness soon after and was bedridden for the rest of my stay. Truly an unfortunate end to an unfortunate trip.",
							"My deadline was approaching and I was forced to abandon the quest - perhaps one day we will learn the truth of this place. "
						}
						Subject3 = "Alas, " .. emotionTable[reporterType][rand(1,#emotionTable[reporterType])] .. " readers, " .. randExcuseTable[rand(1,#randExcuseTable)]
					end
					local section3toFinalBridge = ""
					if numSubjects < 2 then
						local conclusionTable = {
						"Despite a lack of information, I believed I understood the inner truth of this place: ",
						"Although my research was cut short, I knew in my heart just what sort of colony this was: ",
						"Due to my keen intuition, I was able to discern a throughline despite my lack of solid clues: "
						}
						section3toFinalBridge = conclusionTable[rand(1,#conclusionTable)]
					else
						local conclusionTable = {
						"With a wealth of hearsay and anecdote, discerning the truth of this Colony was simple. ",
						"With my research complete, I had a solid grasp on just what sort of Colony this was. ",
						"Finally, I had a full picture of the State of Things. "
						}
						section3toFinalBridge = conclusionTable[rand(1,#conclusionTable)]
					end
					
					local Conclusion = ""
					local conclusionTable = {}
					if columnType == "murder mystery" then
						if positivity >= 7 then --good
							table.insert(conclusionTable,"The culprit was -!! Find out in next week's episode! ")
							table.insert(conclusionTable,"The key was unquestionably the -!! Find out in my tie-in novel, 'The Dregs of Misery'. ")
							table.insert(conclusionTable,"Stay tuned for next week, when the truth is revealed! ")
						elseif positivity <= 3 then --bad
							table.insert(conclusionTable,"This place was naught but a red herring - I needed to leave, posthaste! ")
							table.insert(conclusionTable,"This place was a trick by the true culprit! My investigation had wasted valuable time. ")
							table.insert(conclusionTable,"This colony had nothing to do with the crime at all! A dire turn for our poor investigator. ")
						else --4 5 6
							table.insert(conclusionTable,"This was but a mere stepping-stone on the path to the culprit. ")
							table.insert(conclusionTable,"The key evidence was just out of my grasp - I resolved to travel onward, ever seeking the Truth. ")
							table.insert(conclusionTable,"But that is a story for another time. Stay tuned for next week, when we venture deep into Novorus! ")
						end
					else
						table.insert(conclusionTable,"This place - its " .. simpleInterestTable[columnType][rand(1,#simpleInterestTable[columnType])] .. " - there is only one way to describe them. They are well and truly " .. positivityTable[positivity] .. ". ")
						table.insert(conclusionTable,"In short; this place and its " .. simpleInterestTable[columnType][rand(1,#simpleInterestTable[columnType])] .. " are simply " .. positivityTable[positivity] .. ". ")
						table.insert(conclusionTable,"My time in this colony, experiencing its " .. simpleInterestTable[columnType][rand(1,#simpleInterestTable[columnType])] .. ", could only be described as " .. positivityTable[positivity] .. ". ")
					end
					
					local signOffTable = {
						"With that, " .. emotionTable[reporterType][rand(1,#emotionTable[reporterType])] .. " readers, our time is done. Until my next publication, this is " .. state.name .. ", forever your servant of Journalism.",
						"With that, " .. emotionTable[reporterType][rand(1,#emotionTable[reporterType])] .. " readers, my time with you is done. Until our next meeting, fare well. Signed " .. state.name .. ".",
						"With that, " .. emotionTable[reporterType][rand(1,#emotionTable[reporterType])] .. " readers, my time with you is done. Until next time, this is " .. state.name .. ".",
						"And so, " .. emotionTable[reporterType][rand(1,#emotionTable[reporterType])] .. " readers, this column comes to an end. Join me next fortnight.",
						"And so, " .. emotionTable[reporterType][rand(1,#emotionTable[reporterType])] .. " readers, this column comes to an end. Until next week.",
						"And so, " .. emotionTable[reporterType][rand(1,#emotionTable[reporterType])] .. " readers, this column comes to an end. I shall return to you next third-Wednesday.",
						"In the end, " .. emotionTable[reporterType][rand(1,#emotionTable[reporterType])] .. " readers, we must part. But such things are only temporary - even loss is merely a turning of the tide. Some time, some place, we will meet again. Until then: Thank you.",
					}
					
					Conclusion = conclusionTable[rand(1,#conclusionTable)] .. signOffTable[rand(1,#signOffTable)]
					
					article = article .. Starts[rand(1,#Starts)] .. "\n\n" .. Introduction .. Subject1 .. break1 .. section1to2Bridge .. Subject2 .. break2 .. section2to3Bridge .. Subject3 .. break3 .. section3toFinalBridge .. Conclusion
                         
					local finalReaction = "This article should have a " .. positivityDescTable[positivity] .. " effect on our relations with the Empire and other nations."
					
					send("rendCommandManager",
						"odinRendererFYIMessage",
						iconskin, -- iconskin
						icon, -- icon
						"Article Complete", -- header text
						"The article written about our colony is complete. We've received a copy: \n\n" .. article .. "\n\n" .. finalReaction, -- state.gangName .. " has broken our truce, and intends to resume hostilities!",
						"Left-click to read the column, Right-click to dismiss.", -- tooltip string
						"empireTimes", -- alert type (for stacking)
						"", -- imagename for bg
						"low", -- importance: low / high / critical
						nil, -- state.renderHandle, -- object ID
						6 * 60 * 1000, -- duration in ms
						0, -- "snooze" time if triggered multiple times in rapid succession
						nil)
					
					positivity = (positivity - 5) * 4
                         
                         send( query("gameSession","getSessiongOH", "Empire")[1],
						"changeStanding",positivity,nil)
					send( query("gameSession","getSessiongOH", "Stahlmark")[1],
						"changeStanding",positivity,nil)
					send( query("gameSession","getSessiongOH", "Republique")[1],
						"changeStanding",positivity,nil)
					send( query("gameSession","getSessiongOH", "Novorus")[1],
						"changeStanding",positivity,nil)
					
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
						"Column Writer Flees", -- header text
						state.name .. " is fleeing from this dangerous colony! It is safe to assume that we won't be receiving any articles.", -- text description
						"Left-click to zoom. Right-click to dismiss.", -- action string
						"republicainRogueEnd", -- alert type (for stacking)
						"ui\\eventart\\travellingPoet.png", -- imagename for bg
						"low", -- importance: low / high / critical
						state.rogue.id, -- object ID
						60 * 1000, -- duration in ms
						0, -- snooze
						nil)
					
					send(state.rogue,"addTag","exit_map_run")
				end
				
			elseif path_chosen == "real_nope" then
				-- just leave.
				local icon = "decor_uc_icon"
				local iconskin = "ui\\orderIcons.xml"
				
				send("rendCommandManager",
					"odinRendererStubMessage", -- "odinRendererStubMessage",
					iconskin, -- iconskin
					icon, -- icon
					"Column Writer Leaving", -- header text
					state.name .. " is leaving our colony, the offer made having been denied.", -- text description
					"Left-click to zoom. Right-click to dismiss.", -- action string
					"republicainRogueEnd", -- alert type (for stacking)
					"ui\\eventart\\travellingPoet.png", -- imagename for bg
					"low", -- importance: low / high / critical
					state.rogue.id, -- object ID
					60 * 1000, -- duration in ms
					0, -- snooze
					nil)
                    
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
