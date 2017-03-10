event "fishpeople_collect_eggs"
<<
	state 
	<<
	>>

	receive Create( stringstringMapHandle name )
	<< 
		printl("events", "Fishpeople egg collection event started")
	>>
	
	FSM 
	<<
		["start"] = function(state,tags)
			local targetEggs = tags["fishpeople_eggs"].target
			local targetPos = query(targetEggs,"gridGetPosition")[1]
		
			settimer("Fishpeople Event Timer", 0)
			
			-- Let's add some mystery.
			if rand(1,2) == 1 then
				local s = "unidentified figures"
				local icon = "mysterious_figures"
				local t = "Figures"
				local b = ""
				local fishSeen = query("gameSession", "getSessionBool", "fishpeopleFirstContact")[1]
				if fishSeen then
					-- even more mystery
					if rand(1,2) == 1 then
						s = "Fishpeople"
						icon = "fishperson"
						t = "Fishpeople"
						b = "ui//eventart//fishPeople.png"
					end
				end
				
				send("rendCommandManager", "odinRendererTickerMessage",
					"The Imperial Airship Corps reports that " .. s .. " have been spotted.",
					icon,
					"ui\\thoughtIcons.xml")
					send("rendCommandManager",
						 "odinRendererPlaySoundMessage",
						 "alertNeutral")
				
				send("rendCommandManager",
					"odinRendererStubMessage",
					"ui\\thoughtIcons.xml", -- iconskin
					icon, -- icon
					"" .. t .. " Spotted", -- header text
					"The Imperial Airship Corps reports that " .. s .. " have been spotted.", -- text description
					"Right-click to dismiss.", -- action string
					"fishpeople_eggs", -- alert type (for stacking)
					b, -- imagename for bg
					"low", -- importance: low / high / critical
					nil, -- object ID
					45 * 1000, -- duration in ms
					0, -- "snooze" time if triggered multiple times in rapid succession
					nullHandle)
			end
			
			local spawnTable = { legacyString = "Fishy Patrol Group" }
               local fish_group = query("scriptManager",
								"scriptCreateGameObjectRequest",
								"fishpeople_group",
								spawnTable)[1]
			
               send(fish_group,"GameObjectPlace", -1, -1 )
               send(fish_group,"pushMission","collect_eggs",targetPos, 1 )
			
			return "final"
		end,

		["final"] = function(state,tags)  	
			return
		end,

		["abort"] = function(state,tags) 
			return
		end
	>>
>>
