gameobject "grave"
<<
	local 
	<<
		function epitaph_choose (grave_who_tags, grave_who)
		
			local epitaph_choice = {"REDACTED"}
			
			-- all the epitaphs
			local epitaphs_all = {
				"%s: lover, friend, human... probably.",
				"R.I.P. %s",
				"%s was taken too soon by REDACTED.",
				"Here lies %s.",
				"%s - \"What could go wrong?\"",
				"%s is buried here.",
				"Rest in peace, %s.",
				"%s: Finally free.",
				"%s: Never wanted to be here anyways.",
				"%s: One of the lucky ones.",
				"%s is buried here. The rest of the gravemarker is unintelligible.",
				"%s: REDACTED",
				"Here lies %s. The rest of the gravemarker is covered in claw marks and dried blood.",
				"%s - \"It was worth it!\"",
				"%s - \"It wasn't worth it!\"",
				"%s: What a waste of perfectly good meat..."
			}
			
			local epitaphs_lc = {
				"%s lived and loved, but mostly re-arranged stockpiles.",
				"Here lies %s: loyal peon and meat-enthusiast.",
				"%s: Devoted friend and employee.",
				"%s's job position now available. No health benefits. Minor risk involved.",
				"%s - \"I sure hope there is work after death.\""
			}
			
			local epitaphs_mc = {
				"%s knew how to enjoy the middle-life: tidying, ignoring work condition complaints, and LOTS of alcohol.",
				"%s - \"Work! Faster! PRODUCTIVITY!\"",
				"%s - \"I may be dead, but that's no excuse for slacking off!\"",
				"%s: A true clockworkian genius at paper shuffling.",
				"%s: Respected, feared, drunk."
			}
			
			local epitaphs_uc = {
				"%s is remembered by this monument, as high as their arbitrary standards for fashion and food.",
				"%s - \"Even in death I'm better than all of you combined.\"",
				"%s: Novus Ordo Seclorum",
				"%s - \"Why did I ever come to this filthy backwater?\"",
				"%s lies here. The rest of the gravemarker is covered in complex-looking symbols."
			}
			
			local epitaphs_bandits = {
				"%s: Finally left behind a life of crime.",
				"%s: Was known mostly for their uncanny ability to steal completely worthless commodities and get shot.",
				"%s: Stole my watch. Joke's on them now, the jerk."
			}
			
			local epitaphs_other = {
				"%s: Nobody really knew or understood them. Rest well.",
				"%s: Probably a foreigner of some kind.",
				"%s: Nothing was really known about them."
			}
			
			-- TODO: military specific stuff? more epitaphs? *better* epitaphs?
			
			-- 50/50 chance for exclusive epitaph
			if rand(1,2) == 1 then
				epitaph_choice = epitaphs_all
			else
				if grave_who_tags["bandit"] and not grave_who_tags["citizen"] then
					epitaph_choice = epitaphs_bandits
				elseif grave_who_tags["citizen"] then
					if grave_who_tags["lower_class"] then
						epitaph_choice = epitaphs_lc
					elseif grave_who_tags["middle_class"] then
						epitaph_choice = epitaphs_mc
					elseif grave_who_tags["upper_class"] then
						epitaph_choice = epitaphs_uc
					else
						-- we shouldn't ever hit this, LOG IT
						printl("CECOMMPATCH - Graves: we got a citizen with no class tag (for " .. grave_who .. ") somehow.")
						epitaph_choice = epitaphs_other
					end
				else
					epitaph_choice = epitaphs_other
				end
			end
			
			-- now pick a random one, format the string, and return it
			return string.format(epitaph_choice[rand(1,#epitaph_choice)], grave_who)
		end
	>>

	state
	<<
		gameGridPosition position
		int renderHandle
		string gabionType
		bool addedJob
		int corpseID
	>>

	receive Create(stringstringMapHandle init)
	<<
		local type = init["legacyString"]

		state.position.x = -1
		state.position.y = -1
		state.gabionType = type
		SELF.tags = {"grave"}

		ready()
		sleep()
	>>

	receive GameObjectPlace( int x, int y ) 
	<<	
		state.position.x = x
		state.position.y = y
		state.renderHandle = SELF.id
		
		local models_lc = {
			"graveyardCogWood00.upm",
			"graveyardCogStone00.upm",
			"graveyardCogStone01.upm",
			--"fences/fencePost01.upm",
			--"fences/fencePost02.upm",
		}
		
		local models_mc = {
			"graveyardHeadstone00.upm",
			"graveyardHeadstone01.upm",
			"graveyardHeadstone02.upm",
		}
		
		local models_uc = {
			"graveyardObelisk00.upm",
			"graveyardObelisk01.upm",
		}
		
		local models_military_lc = {
			"fences/whitePicketFencePost01.upm",
		}
		
		local models_military_mc = {
			"memorial00.upm",
		}
		
		local models_other = {
			"fences/rusticFencePost01.upm",
			"fences/rusticFencePost02.upm",
			"fences/rusticFencePost03.upm",
			"fences/basicRailFencePost01.upm",
			"fences/basicRailFencePost02.upm",
			"fences/basicRailFencePost03.upm",
		}
			
		local models = {}
		local grave_who = "Unknown"
		local grave_epitaph = "Details about who is buried here has been forgotten. Or misplaced. Or was unimportant."

		
		local results = query("gameSpatialDictionary",
						  "allObjectsInRadiusWithTagRequest",
						  state.position,
						  1,"corpse",true)[1]
		
		if results then
			grave_who_tags = query(results[1],"getTags")[1]
			grave_who = query(results[1],"getName")[1]
			--printl("CECOMMPATCH: " .. grave_who)
		end	
		
		-- get the epitaph for the gravestone
		grave_epitaph = epitaph_choose(grave_who_tags, grave_who)
		
		if grave_who_tags["citizen"] then
			-- figure out class-specific gravestones and whatnot
			if grave_who_tags["lower_class"] then
				if grave_who_tags["military"] then
					models = models_military_lc
				else
					models = models_lc
				end
			elseif grave_who_tags["middle_class"] then
				if grave_who_tags["military"] then
					models = models_military_mc
				else
					models = models_mc
				end
			elseif grave_who_tags["upper_class"] then
				models = models_uc
			end
		else
			-- non-citizen, give 'em sticks
			models = models_other
		end
		
		local grave_model = models[rand(1,#models)]
		local grave_rotate = rand(-7,7) -- slight rotation for variety. too much is weird looking
		
		-- memorial00.upm has a bad rotation by default, hackish fix
		if grave_model == "memorial00.upm" then
			grave_rotate = grave_rotate + 180
		end
		
		send("rendStaticPropClassHandler",
			"odinRendererCreateStaticPropRequest",
			SELF,
			"models/constructions/" .. grave_model,
			state.position.x,
			state.position.y)

		-- TODO: check with Micah to see if these are, in fact, sane (I would prefer borders.)

		--[[local occupancyMap = 
		"PPPPP\\".. 
		"PpCpP\\"..
		"PPPPP\\"

		local occupancyMapRotate45 = 
		"PPPPP\\".. 
		"PpCpP\\"..
		"PPPPP\\"]]
		
		local occupancyMap = 
			"..@..\\".. 
			".---.\\".. 
			"@-c-@\\"..
			".---.\\"..
			"..@..\\" 
		local occupancyMapRotate45 = 
			"..@..\\".. 
			".---.\\".. 
			"@-c-@\\"..
			".---.\\"..
			"..@..\\"
			
		
		send("gameSpatialDictionary",
               "registerSpatialMapString",
               SELF,
               occupancyMap,
               occupancyMapRotate45,
               true )
         
		send("gameSpatialDictionary",
               "gridAddObjectTo",
               SELF,
               state.position)
		
		send("rendCommandManager",
			"odinRendererHighlightSquare",
			state.position.x,
			state.position.y,
			255, 255, 255, true)
		
		send("rendStaticPropClassHandler",
			"odinRendererMoveStaticPropGridWithHeight",
			state.renderHandle,
			state.position,
			0)
		
		send("rendStaticPropClassHandler",
			"odinRendererOrientStaticProp",
			state.renderHandle,
			0)
		
		-- some slight rotation... must be done LAST because of odinRendererOrientStaticProp
		send("rendStaticPropClassHandler",
			"odinRendererRotateStaticProp",
			state.renderHandle,
			grave_rotate,
			0.25)
		
		-- TODO: get some interesting grave text, randomized, maybe based on the colonist being buried?
          local tooltipTitle = "Grave of " .. grave_who
          local tooltipDescription = grave_epitaph
          send("rendInteractiveObjectClassHandler",
                    "odinRendererBindTooltip",
                    state.renderHandle,
                    "ui//tooltips//groundItemTooltipDetailed.xml",
                    tooltipTitle,
                    tooltipDescription)

	>>

	receive Update()
	<<

	>>

	respond isObstruction()
	<<
		return "obstructionResult", true
	>>


	receive JobCancelledMessage(gameSimJobInstanceHandle job)
	<<
		state.assignment = nil
		state.addedJob = false;
	>>

	respond gridGetPosition()
	<<
		return "reportedPosition", state.position
	>>

	respond gridReportPosition()
	<<
		return "gridReportedPosition", state.position
	>>
	
>>