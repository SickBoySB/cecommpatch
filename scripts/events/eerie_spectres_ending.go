event "eerie_spectres_ending"
<<
	state 
	<<
		bool alertTriggered
          int counter
		string dialogBoxResults
		bool lock
	>>

	receive Create( stringstringMapHandle init )
	<< 
		printl("events","eerie_spectres_ending started")
		state.director = query("gameSession","getSessiongOH","event_director" .. init.director_name)[1]
		state.counter = 4800 -- one day
		state.dialogBoxResults = ""
		state.vicar = query(state.director,"getKeyGameobject","vicar")[1]
	>>
	
	receive boxSelectionResult( string id, string message )
	<<
		state.dialogBoxResults = message
	>>
	
	FSM 
	<<
		["start"] = function(state,tags)
			if state.vicar then
				state.vicarname = query(state.vicar,"getName")[1]
			end
			local ending = query(state.director, "getKeyString", "finale")[1]
			if ending == "success" then
				send("rendCommandManager",
						"odinRendererStubMessage",
						"ui\\thoughtIcons.xml", -- iconskin
						"spectre", -- icon
						"Exorcism success!", -- header text
						state.vicarname .. " has successfully exorcised the angry Spectres! We can all rest peacefully, now.", -- text description
						"Right-click to close.", -- action string
						"spectres", -- alert type (for stacking)
						"ui\\eventart\\eldritch_tinkering.png", -- imagename for bg
						"low", -- importance: low / high / critical
						nil, -- object ID
						60 * 1000, -- duration in ms
						0,
						state.director ) -- "snooze" time if triggered multiple times in rapid succession
                    send("rendCommandManager",
				"odinRendererPlaySoundMessage",
				"alertGood")
				send(state.director, "forwardEventStatusText", "The Spectres were banished by our noble Vicar.", "The Spectres have been banished by our noble Vicar.")
			else
				send("rendCommandManager",
						"odinRendererStubMessage",
						"ui\\thoughtIcons.xml", -- iconskin
						"spectre", -- icon
						"Spectre Trauma!", -- header text
						"Our colony has been thoroughly traumatized by the rampaging spectres.", -- text description
						"Right-click to close.", -- action string
						"spectres", -- alert type (for stacking)
						"ui\\eventart\\eldritch_tinkering.png", -- imagename for bg
						"low", -- importance: low / high / critical
						nil, -- object ID
						60 * 1000, -- duration in ms
						0,
						state.director ) -- "snooze" time if triggered multiple times in rapid succession
                    send("rendCommandManager",
				"odinRendererPlaySoundMessage",
				"alertBad")
				send(state.director, "forwardEventStatusText", "The Spectres traumatized our fair colony.", "The Spectres have traumatized our fair colony.")
			end
			send(state.director,"releaseEventBuilding")
			
			local numSpectres = query(state.director,"getKeyInt","spectreCount")[1]
			if numSpectres then
				if numSpectres > 1 then
					for i = 1,numSpectres do --Kill them Spectres!
						local handle = query(state.director,"getKeyGameobject","spectre" .. i)[1]
						send(handle,"Sunrise")
					end
				end
			end
			
			-- CECOMMPATCH - "Who You Gonna Call" achievement fix
			if ending == "success" then
				local specnum = 1
				
				if numSpectres then
					if numSpectres > 1 then
						specnum = numSpectres
					end
				end
				
				send("gameSession","incSessionInt","spectresBanished", specnum)
					
				local num = query("gameSession","getSessionInt","spectresBanished")[1]
				if not query("gameSession","getSessionBool","whoYouGonnaCall")[1] and num >= 50 then
					send("gameSession", "setSessionBool", "whoYouGonnaCall", true)
					send("gameSession", "setSteamAchievement", "whoYouGonnaCall")
				end
				
				-- There is no way to find out what the necessary steam stat variable is
				-- send("gameSession", "incSteamStat", "stat_whoYouGonnaCall", specnum)
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
