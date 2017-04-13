event "foreign_relations_correction"
<<
	state 
	<<
	>>

	receive Create( stringstringMapHandle init )
	<<
	>>
	
	FSM 
	<<
		["start"] = function(state,tags)
			printl("events", "foreign_relations_correction started")
			
			local r = rand(1,4)
			local nationNameShort = false
			local targetRelations = 0
			
			if r == 1 then
				-- correct allied
				nationNameShort = query("gameSession","getSessionString","defaultFriendly")[1]
				targetRelations = 75
				
			elseif r == 2 then
				-- correct neutral
				nationNameShort = query("gameSession","getSessionString","defaultNeutral")[1]
				
			elseif r == 3 then
				-- correct hostile
				nationNameShort = query("gameSession","getSessionString","defaultHostile")[1]
				targetRelations = -75
				
			else
				-- do dominant foreign pressence
				local dominant = query("gameSession","getSessionString", "dominantFaction")[1]
				local hostile = query("gameSession","getSessionString","defaultHostile")[1]
				local neutral = query("gameSession","getSessionString","defaultNeutral")[1]
				local friendly = query("gameSession","getSessionString","defaultFriendly")[1]
				
				-- CECOMMPATCH more effing hacks for dominant not using the right name
				local dominant_short_name = {
					["Grossherzoginnentum von Stahlmark"] = "Stahlmark",
					["Novorus Imperiya"] = "Novorus",
					["Republique Mecanique"] = "Republique",
					-- not sure if any of the below is ever used, but just to be safe..
					["The Clockwork Empire"] = "Empire",
					["Fishpeople"] = "Fishpeople",
					["Bandits"] = "Bandits"
				}
					
				if dominant_short_name[dominant] then
					printl("CECOMMPATCH - dominant faction shortname correction")
					dominant = dominant_short_name[dominant]
				end
				-- /hack
				
				if hostile == dominant then
					nationNameShort = dominant
					targetRelations = -75
				elseif neutral == dominant then
					nationNameShort = dominant
				elseif friendly == dominant then
					nationNameShort = dominant
					targetRelations = 75
				end
			end
			
			local currentRelations = query("gameSession", "getSessionInt", nationNameShort .. "Relations")[1]

			local relationsChange = targetRelations - currentRelations
			
			if relationsChange == 0 then
				printl("events", "foreign_relations_correction: no change needed, ending.")
				return "final"
			end
			printl("events", "foreign_relations_correction: current/target = " .. currentRelations .. " / " .. targetRelations )
			
			state.nation = query("gameSession","getSessiongOH", nationNameShort )[1]
			local nationInfo = EntityDB[  nationNameShort .. "Info" ]
			local nationName = nationInfo.fullName
			local theirLeadership = nationInfo.leadershipRepresentativeString
			
			local increase = false
			if relationsChange > 0 then
				increase = true
			end
			
			local currentState = query(state.nation,"getRelationStateString")[1]
			
			local dialogString = "" 
			local dialogTitle = "Foreign Relations Change"
			local artPath = "ui\\eventart\\capitalists.png"
			-- oh god here we go
			
			if currentState == "hostile" then
				if increase then
					local r = rand(1,2)
					if r == 1 then
						dialogString = "Although simmering in a low-stakes war, diplomats have arranged with " .. theirLeadership .. " for an exchange of prisoners and merchant vessels caught in the conflict.\n\nThis improves relations slightly between " .. nationName .. " and the Empire, although hostilities continue."	
						dialogTitle = "Prisoner Exchange" -- with the " .. nationName
					elseif r == 2 then
						dialogString =  "Representatives of Her Majesty, on behalf of the Clockwork Empire, have met with the Foreign Minister of " .. nationName .. " at a conference held in the neutral territory of the Alpine Confederation.\n\nAlthough inconclusive, much progress was made in defining the terms of the conflict, and relations have improved as a result."
						dialogTitle = "Diplomatic Conference" -- with the " .. nationName
					end
				else -- decrease
					local r = rand(1,2)
					if r == 1 then 
						local battleTypes = {
							"battle in the skies above the Great Channel between mighty airships, raining bodies and shattered warcraft alike",
							"terrible exchange of battery-fire and torpedo between great warships of - and under - the sea, with many hundreds feared drowned",
							"series of colonial border skirmishes resulting in commitment of several squads of our noble Steam Knights to counter their unprovoked advance",
							}
						local battleType = battleTypes[ rand(1,#battleTypes) ]
						
						dialogString = "The " .. nationName .. " has fought the Empire in a " .. battleType .. "! This has caused much anger amongst the common folk and betters of the Empire alike toward " .. nationName .. " due to their barbarous offenses. \n\nRelations are decreased as a result."
						dialogTitle = "News Of Battle" -- with the " .. nationName
						artPath = "ui\\eventart\\battle.png"
					else
						dialogString = "The pandering and loathesome Foreign Minister of " .. nationName .. " has publically insulted the Clockwork Empire at a diplomatic conference being held on the neutral ground of the Alpine Confederation.\n\nThis has damaged relations between us."
						dialogTitle = "Diplomatic Insult" -- from the " .. nationName
					end
					
				end
			elseif currentState == "friendly" then
				
				if increase then
					dialogString ="Representatives of " .. theirLeadership .. " were invited to tour the Capital, charming and being charmed in return by Clockworkian aristocrats, ministers, and Barons of Industry.\n\nRelations between the Empire and " .. nationName .. " have improved as a result!"
					dialogTitle = "Diplomatic Tour" -- from the " .. nationName
				else -- decrease
					
					local actions = {
							"seizing several Clockworkian merchant airships amidst accusations of smuggling and piracy.",
							"burning down a number of key heliograph repeaters on the distant frontier while accusing the Empire of infringing upon claimed territory.",
							"declaring that Clockworkian goods are inferior and their importation constitutes damages against their State.",
						}
					
					local hostileActionDescription = actions[ rand(1,#actions) ] 
						
					dialogString = "Although the colonial forces of " .. nationName .. " are well-disposed toward the Empire in general and you in particular, " ..  theirLeadership .. " have pushed a harder line against the Clockwork Empire, " .. hostileActionDescription .. " \n\nOur relations with " .. nationName .." have degraded as a result."
					dialogTitle = "Hardline Tactics" -- from the " .. nationName
				end
					
			else -- neutral
				if increase then
					dialogString =  "Representatives of " .. theirLeadership .. " were invited to tour the Capital, charming and being charmed in return by Clockworkian aristocrats, ministers, and Barons of Industry.\n\nRelations between the Empire and " .. nationName .. " have improved as a result!"
					
				else -- decrease
					local r = rand(1,2)
					if r == 1 then 
						dialogString = "The pandering and loathesome Foreign Minister of " .. nationName .. " has publically insulted the Clockwork Empire at a diplomatic conference being held on the neutral ground of the Alpine Confederation.\n\nThis has damaged relations between us."
						dialogTitle = "Diplomatic Insult" -- from the " .. nationName
					else
						local actions = {
							"seizing several Clockworkian merchant airships amidst accusations of smuggling and piracy.",
							"burning down a number of key heliograph repeaters on the distant frontier while accusing the Empire of infringing upon claimed territory",
							"declaring that Clockworkian goods are inferior and their importation constitutes damages against their State.",
						}
						local hostileActionDescription = actions[ rand(1,#actions) ]
						
						dialogString = "Although the colonial bureacracy of " .. nationName .. " is neutrally disposed toward your colony in particular, orders have come from " .. theirLeadership .. " to push a hard line against the Clockwork Empire, resulting in them " .. hostileActionDescription .. " \n\nOur relations with " .. nationName .." have degraded as a result."
						dialogTitle = "Hardline Tactics" -- from the " .. nationName
					end
				end
			end
			
			send("rendCommandManager",
				"odinRendererFYIMessage",
				nationInfo.iconSkin,
				nationInfo.iconFlag,
				dialogTitle,
				dialogString,
				"Left-click for details.",
				dialogTitle,
				artPath,
				"high",
				nil,
				60 * 1000,
				0, -- "snooze" time if triggered multiple times in rapid succession
				nil) -- gameobjecthandle of director, null if none
			
			-- nothing too extreme, please.
			relationsChange = math.floor( relationsChange * 0.5 )
			
			if relationsChange > 15 then relationsChange = 15 end
			if relationsChange < -10 then relationsChange = -10 end
			
			send(state.nation,"changeStanding",relationsChange,nil)
			
			return "final"	
		end,

		["final"] = function(state,tags)
			settimer("Foreign Event Timer",0)
		
			return
		end,

		["abort"] = function(state,tags) 
			return
		end
	>>
>>