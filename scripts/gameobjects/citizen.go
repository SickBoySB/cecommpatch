gameobject "citizen" inherit "ai_agent"
<<

	local
	<<
		function makeMemory(memoryName, memoryDescription, otherName, otherObject, otherObjectKey)
			
			if SELF.tags.selenian_infested then
				memoryName = "The Stars Are Right"
				memoryDescription = nil
				otherName = nil
				otherObject = nil
				otherObjectKey = nil
			end
			
			local memory = EntitiesByType["emotion"][memoryName]
			if not memory then
				printl("ai_agent", state.AI.name .. " : warning, couldn't find memory! " .. memoryName)
				return
			end
			
			if not memoryDescription then
				memoryDescription = memory.info.description
			end

			memoryDescription = parseDescription(memoryDescription)
			if otherName then
				memoryDescription = string.gsub(memoryDescription, "OTHER_NAME", otherName )
			end
			
			local characters_involved = { character = SELF }
			
			if otherObject and otherObjectKey then
				characters_involved[otherObjectKey] = otherObject
			end
			
			send("gameHistoryDB","createHistoryFragment", 
				memory.name, 
				memory.type,
				characters_involved,		
				{	description = memoryDescription,
					painting_description = memory.info.painting_description,
					priority = memory.info.priority,
					duration = memory.info.duration, 
					icon = memory.info.icon,
					iconSkin = memory.info.iconSkin }, 
				memory.values )
		end
		
		function makeWorkCrewNames()
			-- set up name for civ and mil versions of workcrew.
				
			local stats = EntityDB.HumanStats
			
			local prefix = ""
			local modifier = ""
			
			state.AI.strs.workCrewNameCivilian = ""
			state.AI.strs.workCrewNameMilitary = ""
			
			local iter = 1
			while iter == 1 or string.len(state.AI.strs.workCrewNameCivilian) > 38 do
				local baseCivilianName = "Workcrew"
				local highestSkill = false
				local highestSkillAmount = 0
				for skillName,skillAmount in pairs(state.AI.skills) do
					if skillAmount > highestSkillAmount then
						highestSkillAmount = skillAmount
						highestSkill = skillName
					end
				end
				
				if highestSkillAmount > 1 and
					stats.workcrewNamesBySkill[highestSkill] then
					
					-- do skill-based naming.
					local nameSet = stats.workcrewNamesBySkill[highestSkill]
					baseCivilianName = nameSet[ rand(1,#nameSet) ]
				else
					-- do generic name, maybe based on assignment.
					if state.AI.strs["citizenClass"] == "Barber" then
						baseCivilianName = "Crude Medicine and Haircuts"
					elseif state.AI.strs["citizenClass"] == "Vicar" then
						baseCivilianName = "Congregation"
                         elseif state.AI.strs["citizenClass"] == "Trainee" then
						baseCivilianName = "Class"    
					elseif state.AI.strs["citizenClass"] == "Naturalist" then
						baseCivilianName = "Expedition"
					elseif state.AI.strs["citizenClass"] == "Artisan" then
						baseCivilianName = "Artisanal Company"
					-- Do one for prisoners?
					elseif state.AI.traits["Prison Overseer"] == true then
						baseCivilianName = "Convict Labourers"
					else
						local names = EntityDB.HumanStats.workcrewNamesBasic
						baseCivilianName = names[ rand(1,#names)]
					end
				end
				
				-- now get a trait descriptor
				for traitName,v in pairs(state.AI.traits) do
					if v and stats.workcrewNameModsByTrait[traitName] then
						local set = stats.workcrewNameModsByTrait[traitName]
						modifier = set[ rand(1,#set) ]
						break
					end
				end
				
				state.AI.strs.workCrewNameCivilian = state.AI.name .. "'s " .. modifier .. " " .. baseCivilianName
				iter = iter + 1
				if iter > 10 then
					break
				end
			end
			
			local iter = 1
			while iter == 1 or string.len(state.AI.strs.workCrewNameMilitary) > 38 do
				-- Military naming.
				local baseMilitaryName = "Squad"
				local diceroll = rand(10,200)
				local numberEnding = "th"
				if diceroll == 11 or
					diceroll == 111 then
					numberEnding = "th"
				elseif diceroll == 12 or
					diceroll == 112 then
					numberEnding = "th"
				elseif diceroll%10 == 1 then
					numberEnding = "st"
				elseif diceroll%10 == 2 then
					numberEnding = "nd"
				elseif diceroll%10 == 3 then
					numberEnding = "rd"
				end
				
				local unitTypes = stats.workcrewNamesMilitary
				prefix = "Her Majesty's"
				baseMilitaryName = diceroll .. numberEnding .." " .. unitTypes[ rand(1,#unitTypes) ]
				
				-- now get a trait descriptor
				for traitName,v in pairs(state.AI.traits) do
					if v and stats.workcrewNameModsByTrait[traitName] then
						local set = stats.workcrewNameModsByTrait[traitName]
						modifier = set[ rand(1,#set) ]
						break
					end
				end
				
				state.AI.strs.workCrewNameMilitary = prefix .. " " .. modifier .. " " .. baseMilitaryName
				iter = iter + 1
				if iter > 10 then
					break
				end
			end

			return
		end

		function RunDecisionTree ( job )
			local results = query ("gameBlackboard", "gameAgentEvaluateDecisionTreeMessage", state.AI, SELF, job )
			if results.name == "gameAgentAssignedJobMessage" then
				results[1].assignedCitizen = SELF
				if state.AI.curJobInstance then
					if VALUE_STORE["showFSMDebugConsole"] then
						printl("FSM", "FSM: " .. state.AI.strs["firstName"] .. " " .. state.AI.strs["lastName"] .. " is attempting to abort job " .. state.AI.curJobInstance.name .. " due to a decision tree (" .. results[1].name .. ") !" ) 
					end
					if state.AI.FSMindex > 0 then
						-- run the abort state
						local tag = state.AI.curJobInstance:findTag( state.AI.curJobInstance:FSMAtIndex( state.AI.FSMindex  ) )
						local name = state.AI.curJobInstance:findName ( state.AI.curJobInstance:FSMAtIndex( state.AI.FSMindex ) )
						-- load up our fsm 
						local FSMRef = state.AI.curJobInstance:FSMAtIndex(state.AI.FSMindex)
						
						local targetFSM
						if FSMRef:isFSMDisabled() then
							targetFSM = ErrorFSMs[ FSMRef:getErrorFSM() ]
						else 
							targetFSM = FSMs[ FSMRef:getFSM() ]
						end

						local ok
						local nextState
						ok, errorState = pcall( function() targetFSM[ "abort" ](state, tag, name) end )
	
						if not ok then 
							print("ERROR: " .. errorState )
							FSM.stateError( state )
						end
					end

					state.AI.curJobInstance:abort( "Interrupt hit." )
					state.AI.curJobInstance = nil
				end

				state.AI.abortJob = true
				if reason ~= nil and reason ~= "" then
					state.AI.abortJobMessage = reason 
				end 
				if state.AI.abortJobMessage == nil then
					state.AI.abortJobMessage = ""
				end

				-- Reset the counter for next time
				state.AI.FSMindex = 0
				state.AI.curJobInstance = results[ 1 ]
				send("rendOdinCharacterClassHandler",
					"odinRendererSetCharacterAttributeMessage",
						state.renderHandle,
						"currentJob",
						state.AI.curJobInstance.displayName)
				
				send("rendOdinCharacterClassHandler",
					"odinRendererSetCharacterAttributeMessage",
						state.renderHandle,
						"currentJobCategory",
						state.AI.curJobInstance.filter)
				
				return true		
			elseif results.name == "gameAgentDecisionTreeOKMessage" then
				return true -- I am okay. Don't need to set up a new job; but you are doing an okay job in the tree right now.
			end			-- results
			return false
		end	

		function ForceJob ( job, table )

			local results = query( "gameBlackboard", "gameAgentForceJobMessage", state.AI, SELF, job, table )
			if results.name == "gameAgentAssignedJobMessage" then
				results[1].assignedCitizen = SELF
				if state.AI.curJobInstance then
					
					if VALUE_STORE["showFSMDebugConsole"] then
						printl("FSM", "FSM: " .. state.AI.name .. " is attempting to abort job " .. state.AI.curJobInstance.name .. " due to a forced job (" .. results[1].name .. ") !" )
					end
				--[[
					if SELF.tags["military"] then
						local tickerText = state.AI.strs["firstName"] .. " " .. state.AI.strs["lastName"] .. " interrupts job " .. state.AI.curJobInstance.name .. " to do forced job " .. results[1].name
						send("rendCommandManager", "odinRendererTickerMessage", tickerText, "work", "ui\\thoughtIcons.xml")
					end
					--]]
					if state.AI.FSMindex > 0 then
						-- run the abort state
						local tag = state.AI.curJobInstance:findTag( state.AI.curJobInstance:FSMAtIndex( state.AI.FSMindex  ) )
						local name = state.AI.curJobInstance:findName ( state.AI.curJobInstance:FSMAtIndex( state.AI.FSMindex ) )
						
						if not name then name = "nil name" end
						-- load up our fsm 
						local FSMRef = state.AI.curJobInstance:FSMAtIndex(state.AI.FSMindex)
						
						local targetFSM
						if FSMRef:isFSMDisabled() then
							targetFSM = ErrorFSMs[ FSMRef:getErrorFSM() ]
						else 
							targetFSM = FSMs[ FSMRef:getFSM() ]
						end

						local ok
						local nextState
						ok, errorState = pcall( function() targetFSM[ "abort" ](state, tag, name) end )
	
						if not ok then 
							print("ERROR: " .. errorState )
							FSM.stateError( state )
						end
					end

					state.AI.curJobInstance:abort( "Interrupt hit." )
					state.AI.curJobInstance = nil
				end

				state.AI.abortJob = true
				if reason ~= nil and reason ~= "" then
					state.AI.abortJobMessage = reason 
				end 
				if state.AI.abortJobMessage == nil then
					state.AI.abortJobMessage = ""
				end

				-- Reset the counter for next time
				state.AI.FSMindex = 0
				state.AI.curJobInstance = results[ 1 ]
				send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage",
						state.renderHandle,
						"currentJob",
						state.AI.curJobInstance.displayName)
				send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage",
						state.renderHandle,
						"currentJobCategory",
						state.AI.curJobInstance.filter)

				end			-- if
		end			-- function
	
	function tier1JobDecisionTree()
		-- printl("Job test", state.AI.strs["firstName"] .. state.AI.strs["lastName"] .. " STARTED running tier1JobDecisionTree")
		-- TODO: irrelevant until implementation of fire.
		--[[if state.AI.bools["onFire"] then
			if state.AI.traits["Doomed"] then
				ForceJob("Idle (Human)", {})
				return true
			else
				ForceJob("Flee (Human)",{})
				return true
			end
		end--]]
		return false
	end

	function tier2JobDecisionTree()
		if SELF.tags.frontier_justice and not
			SELF.tags.doing_murder_rampage and not
			SELF.tags.cannibalistic_murderer and not
			SELF.tags.cannibal then
			
			ForceJob("Flee Far Away (Human)",{})
		
			return true
		end
		return false
	end
	
	function tier3MilitaryJobDecisionTree()
		--Put check for moving out of ignition range of fire here.
		-- printl("Job test", state.AI.strs["firstName"] .. state.AI.strs["lastName"] .. " FINISHED running tier3MilitaryJobDecisionTree")
		return false
	end

	function tier3CivilianJobDecisionTree()
		--Put check for moving out of ignition range of fire here.
		return false
	end

	function tier4MilitaryJobDecisionTree()

		-- just run the entire decision tree.
		if SELF.tags["military_decisiontree"] == true then			
			if RunDecisionTree("Military Tree: Raise Alarm") == true then
				return true
			end
		end

		return false

	end

	function tier4CivilianJobDecisionTree()
		
		if RunDecisionTree("Civilian Tree: Raise Alarm") == true then
			return true
		end

		return false
	end

	function tier5MilitaryJobDecisionTree()

	--[[
		if state.AI.curJobInstance and state.AI.curJobInstance.name == "Respond To Alarm" then
			return true
		end
		if state.AI.curJobInstance and state.AI.curJobInstance.name == "Reload Firearm" then
			return true
		end

		-- use reverse lookup so we don't open up 10,000 nodes
		local alarmResult = query("gameSpatialDictionary",
							 "isObjectInRadiusWithTagReverse",
							 state.AI.position,
							 50,
							 "alarm_waypoint")
		
		if alarmResult and alarmResult[1] then			
			ForceJob("Respond To Alarm", {})
			if state.AI.curJobInstance and state.AI.curJobInstance.name == "Respond To Alarm" then
				return true
			end
		end
		--]]
		return false
	end

	
	function lowUrgencySurvivalDecisionTree()
		-- food & sleep
		
		-- This is all waiting on decision tree overhaul. -dgb
		--[[if state.AI.ints.hunger >= 1 then
			-- eat cooked/pickled food
			local rangeResults = query("gameSpatialDictionary",
								    "isObjectInRadiusWithTagReverseMustOwn",
								    state.AI.position, 100,
								    "cooked", "food")
			
			if rangeResults[1] then
				if state.AI.bools["isWieldingTool"] then
					local item_tags = query( state.AI.possessedObjects["curCarriedTool"], "getTags" )[1]
					if item_tags.tool then
						ForceJob("Desummon Tool", {})
					elseif item_tags.weapon then
						ForceJob("Stow Tool", {})
					end
				end
				ForceJob("Eat Cooked Food", {})
				return true
			end
		end
		
		if state.AI.ints.hunger >= 2 then
			-- begrudgingly eat raw food
			local rangeResults = query("gameSpatialDictionary",
								    "isObjectInRadiusWithTagReverseMustOwn",
								    state.AI.position, 100,
								    "raw", "food")
			
			if rangeResults[1] then
				if state.AI.bools["isWieldingTool"] then
					local item_tags = query( state.AI.possessedObjects["curCarriedTool"], "getTags" )[1]
					if item_tags.tool then
						ForceJob("Desummon Tool", {})
					elseif item_tags.weapon then
						ForceJob("Stow Tool", {})
					end
				end
				ForceJob("Eat Raw Food", {})
				return true
			end
		end
		
		if state.AI.ints.hunger >= 3 then
			-- forage due to hunger
			local rangeResults = query("gameSpatialDictionary",
								    "isObjectInRadiusWithTagReverse",
								    state.AI.position, 100,
								    "forage_food_source")
			
			if rangeResults[1] then
				if state.AI.bools["isWieldingTool"] then
					local item_tags = query( state.AI.possessedObjects["curCarriedTool"], "getTags" )[1]
					if item_tags.tool then
						ForceJob("Desummon Tool", {})
					elseif item_tags.weapon then
						ForceJob("Stow Tool", {})
					end
				end
				
				ForceJob("Forage Due To Hunger", {})
				return true
			end
		end
		
		if state.AI.ints.hunger >= 4 then
			-- hunt animals by hand & eat raw meat
			-- & butcher fishperson corpses for meat
			local rangeResults = query("gameSpatialDictionary",
								    "isObjectInRadiusWithTagReverse",
								    state.AI.position, 100,
								    "dead_fishperson")
			
			if rangeResults[1] then
				if state.AI.bools["isWieldingTool"] then
					local item_tags = query( state.AI.possessedObjects["curCarriedTool"], "getTags" )[1]
					if item_tags.tool then
						ForceJob("Desummon Tool", {})
					elseif item_tags.weapon then
						ForceJob("Stow Tool", {})
					end
				end
				
				ForceJob("Butcher Fishperson Corpse", {})
				return true
			end
		end
		
		if state.AI.ints.hunger >= 5 then
			-- butcher human corpses for meat
			local rangeResults = query("gameSpatialDictionary",
								    "isObjectInRadiusWithTag",
								    state.AI.position, 30,
								    "dead_human")
			
			if rangeResults[1] then
				if state.AI.bools["isWieldingTool"] then
					local item_tags = query( state.AI.possessedObjects["curCarriedTool"], "getTags" )[1]
					if item_tags.tool then
						ForceJob("Desummon Tool", {})
					elseif item_tags.weapon then
						ForceJob("Stow Tool", {})
					end
				end
				
				ForceJob("Butcher Human Corpse", {})
				return true
			end
		end
		
		if state.AI.ints.hunger == 6 then
			-- if traits/despair, hunt humans for meat (handled elsewhere)
		end
		
		if state.AI.ints.hunger == 7 then
			-- death due to starvation. (handled elsewhere.)
		end--]]
		
		
		return false
	end
	
	-- This update occurs once per three game seconds.
	function character_doThreeSecondUpdate()
		tmEnter("character three second update")
		send("gameDesiresHandler", "EvaluateDesires", SELF, state.AI)

		if state.AI.ints["inebriation"] > 0 then
			state.AI.ints["inebriation_counter"] = state.AI.ints["inebriation_counter"] + 1
			if state.AI.ints["inebriation_counter"] >= state.AI.ints["3s_per_inebriation"] then
				state.AI.ints["inebriation"] = state.AI.ints["inebriation"] - 1
				state.AI.ints["inebriation_counter"] = 0
			end
		end
		
		if SELF.tags.lower_class and
			not SELF.tags.dead then
			
			local overseer = query("gameBlackboard",
						   "gameObjectGetOverseerMessage",
						   state.AI.currentWorkParty)[1]
				
			
		end
			
		mentalStateAggregator()
		makeMoodText()
		setDescriptiveParagraph()

		local parsedParagraph = parseDescription(state.descriptiveParagraph)
		send("rendOdinCharacterClassHandler",
				"odinRendererSetDescriptionParagraph",
				state.renderHandle,
				parsedParagraph)
		tmLeave()

	end

	-- This update occurs once per game second.
	function character_doSecondUpdate()
		
		tmEnter("character_doSecondUpdate")
		state.AI.ints["alarmTimer"] = state.AI.ints["alarmTimer"] + 1
		if state.AI.ints["alarmTimer"] > 10 then
			if SELF.tags["alarm_waypoint_active"] then
				SELF.tags["alarm_waypoint_active"] = nil
				send("gameSpatialDictionary", "gameObjectRemoveBit", SELF, 17);
			end
		end
		
		state.AI.ints["emotionAnimationTimer"] = state.AI.ints["emotionAnimationTimer"] + 1
		
		-- health regen
		if state.AI.ints["health"] < 10 then
			state.AI.ints["healthRegenTimer"] = state.AI.ints["healthRegenTimer"] - 1
			
			if state.AI.ints["healthRegenTimer"] <= 0 then
				state.AI.ints["health"] = state.AI.ints["health"]+1

				state.AI.ints["healthRegenTimer"] = EntityDB["HumanStats"]["healthRegenTimerSeconds"]
			end
		end

		if state.myAfflictions and #state.myAfflictions > 0 then
			SELF.tags.injured = true
		else
			SELF.tags.injured = false
		end			

		
		mentalStateAggregator()

		--Supply counter and tracking OC-2385

		if state.AI.ints["inCombatTimer"] > 0 then
			state.AI.ints["inCombatTimer"] = state.AI.ints["inCombatTimer"] - 1
			if state.AI.ints["inCombatTimer"] == 0 then
				send("rendOdinCharacterClassHandler", "removeCombatPanel", SELF.id);
			end
		end

		local parsedParagraph = parseDescription(state.descriptiveParagraph)
		send("rendOdinCharacterClassHandler",
				"odinRendererSetDescriptionParagraph",
				state.renderHandle,
				parsedParagraph)
		
		tmLeave()
	end

	function parseDescription(descriptionParagraph)
		
		if string.find(descriptionParagraph, 'CHARACTER_NAME') then
			descriptionParagraph = string.gsub(descriptionParagraph, "CHARACTER_NAME", state.AI.strs["firstName"] .. " " .. state.AI.strs["lastName"])
		end
		
		if string.find(descriptionParagraph, 'CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN') then
			if state.AI.strs["genderstr"] == "male" then
				descriptionParagraph = string.gsub(descriptionParagraph, "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN", "He")
			elseif state.AI.strs["genderstr"] == "female" then
				descriptionParagraph = string.gsub(descriptionParagraph, "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN", "She")
			else 
				descriptionParagraph = string.gsub(descriptionParagraph, "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN", "They")
			end
		end
		
		if string.find(descriptionParagraph, 'SUBJECTIVE_PERSONAL_PRONOUN') then
			if state.AI.strs["genderstr"] == "male" then
				descriptionParagraph = string.gsub(descriptionParagraph, "SUBJECTIVE_PERSONAL_PRONOUN", "he")
			elseif state.AI.strs["genderstr"] == "female" then
				descriptionParagraph = string.gsub(descriptionParagraph, "SUBJECTIVE_PERSONAL_PRONOUN", "she")
			else 
				descriptionParagraph = string.gsub(descriptionParagraph, "SUBJECTIVE_PERSONAL_PRONOUN", "they")
			end
		end
		
		if string.find(descriptionParagraph, 'CAPITALIZED_POSSESSIVE_PERSONAL_PRONOUN') then
			if state.AI.strs["genderstr"] == "male" then
				descriptionParagraph = string.gsub(descriptionParagraph, "CAPITALIZED_POSSESSIVE_PERSONAL_PRONOUN", "His")
			elseif state.AI.strs["genderstr"] == "female" then
				descriptionParagraph = string.gsub(descriptionParagraph, "CAPITALIZED_POSSESSIVE_PERSONAL_PRONOUN", "Her")
			else 
				descriptionParagraph = string.gsub(descriptionParagraph, "CAPITALIZED_POSSESSIVE_PERSONAL_PRONOUN", "Their")
			end
		end
		
		if string.find(descriptionParagraph, 'POSSESSIVE_PERSONAL_PRONOUN') then
			if state.AI.strs["genderstr"] == "male" then
				descriptionParagraph = string.gsub(descriptionParagraph, "POSSESSIVE_PERSONAL_PRONOUN", "his")
			elseif state.AI.strs["genderstr"] == "female" then
				descriptionParagraph = string.gsub(descriptionParagraph, "POSSESSIVE_PERSONAL_PRONOUN", "her")
			else 
				descriptionParagraph = string.gsub(descriptionParagraph, "POSSESSIVE_PERSONAL_PRONOUN", "their")
			end
		end
		
		if string.find(descriptionParagraph, 'OBJECTIVE_PERSONAL_PRONOUN') then
			if state.AI.strs["genderstr"] == "male" then
				descriptionParagraph = string.gsub(descriptionParagraph, "OBJECTIVE_PERSONAL_PRONOUN", "him")
			elseif state.AI.strs["genderstr"] == "female" then
				descriptionParagraph = string.gsub(descriptionParagraph, "OBJECTIVE_PERSONAL_PRONOUN", "her")
			else 
				descriptionParagraph = string.gsub(descriptionParagraph, "OBJECTIVE_PERSONAL_PRONOUN", "them")
			end
		end

		return descriptionParagraph
	end

	function parseMemoriesForEmotions()
		send("rendOdinCharacterClassHandler",
				"odinRendererSetCharacterAttributeMessage",
				state.renderHandle,
				"mood",
				state.AI.strs["mood"] )

		send(SELF,"refreshCharacterAlert")
		
		if state.AI.strs["mood"] ~= state.AI.strs["oldMood"] then
			send("rendCommandManager", "odinRendererRegenerateColonistTableMessage")
		end

		setDescriptiveParagraph()
	end

	function emote()
		if state.AI.strs["socialClass"] == "lower" then
			if state.AI.ints["emoteTimer"] > 45 then
				send(SELF, "emoteThought")
				send(SELF, "resetEmoteTimer")
			end
		else 
			if state.AI.ints["emoteTimer"] > 20 then
				send(SELF, "emoteThought")
				send(SELF, "resetEmoteTimer")
			end
		end
	end


	function mentalStateAggregator()
		
		state.AI.strs["oldMood"] = state.AI.strs["mood"]
		local resultMemories = query("gameHistoryDB",
			"UpdateMoodFromHistory",
			SELF,
			state.AI)

		parseMemoriesForEmotions()

		local i = 1
		
		-- need to do this because when we send a nil value through to C++ it seems to set it to zero

		for k=1, #state.AI.longTermMemories do
			local thisMemory = state.AI.longTermMemories[k]
			if thisMemory.attributes.happiness == nil then thisMemory.attributes.happiness = 0 end
			if thisMemory.attributes.despair == nil then thisMemory.attributes.despair = 0 end
			if thisMemory.attributes.fear == nil then thisMemory.attributes.fear = 0 end
			if thisMemory.attributes.anger == nil then thisMemory.attributes.anger = 0 end
		end
		if state.AI.shortTermMemory.attributes.happiness == nil then state.AI.shortTermMemory.attributes.happiness = 0 end
		if state.AI.shortTermMemory.attributes.despair == nil then state.AI.shortTermMemory.attributes.despair = 0 end
		if state.AI.shortTermMemory.attributes.fear == nil then state.AI.shortTermMemory.attributes.fear = 0 end
		if state.AI.shortTermMemory.attributes.anger == nil then state.AI.shortTermMemory.attributes.anger = 0 end

		for k=1, #state.AI.longTermMemories do
			local thisMemory = state.AI.longTermMemories[k]
			if (state.personalMemories[i] == nil) or (state.personalMemories[i] ~= thisMemory) then
				send("rendCommandManager", "odinRendererDeleteMemoryMessage", SELF, k);
				if thisMemory.data.iconSkin ~= nil then
					send("rendCommandManager", "odinRendererNewMemoryMessage", 
							SELF, i, thisMemory.eventName, thisMemory.data.icon, thisMemory.data.iconSkin, 
							thisMemory.attributes.happiness, thisMemory.attributes.despair, thisMemory.attributes.fear, thisMemory.attributes.anger,
							thisMemory.attributes.importance, thisMemory.data.priority, false, parseDescription(thisMemory.data.description))
				else
					send("rendCommandManager", "odinRendererNewMemoryMessage", 
							SELF, i,  thisMemory.eventName, thisMemory.data.icon, "thoughtIcons.xml", 
							thisMemory.attributes.happiness, thisMemory.attributes.despair, thisMemory.attributes.fear, thisMemory.attributes.anger,
							thisMemory.attributes.importance, thisMemory.data.priority, false, parseDescription(thisMemory.data.description))
				end
				state.personalMemories[ i ] = thisMemory
			end
			i = i+1
		end

		-- same as above, but specifically for the one short term memory, so the odinRendererNewMemoryMessage gets a "true" instead.

		local thisMemory = state.AI.shortTermMemory
		if (state.personalMemories[i] == nil) or (state.personalMemories[i] ~= thisMemory) then
			send("rendCommandManager", "odinRendererDeleteMemoryMessage", SELF, i);
			if thisMemory.data.iconSkin ~= nil then
				send("rendCommandManager", "odinRendererNewMemoryMessage", 
						SELF, i, thisMemory.eventName, thisMemory.data.icon, thisMemory.data.iconSkin, 
						thisMemory.attributes.happiness, thisMemory.attributes.despair, thisMemory.attributes.fear, thisMemory.attributes.anger,
						thisMemory.attributes.importance, thisMemory.data.priority, true, parseDescription(thisMemory.data.description))
			else
				send("rendCommandManager", "odinRendererNewMemoryMessage", 
						SELF, i, thisMemory.eventName, thisMemory.data.icon, "thoughtIcons.xml", 
						thisMemory.attributes.happiness, thisMemory.attributes.despair, thisMemory.attributes.fear, thisMemory.attributes.anger,
						thisMemory.attributes.importance, thisMemory.data.priority, true, parseDescription(thisMemory.data.description))
			end
			state.personalMemories[ i ] = thisMemory
		end
		i = i+1

		if i < 10 then
			for k=i+1,10 do
				send("rendCommandManager", "odinRendererDeleteMemoryMessage", SELF, k)
			end
		end

--[[

		local i = 1
		for k,v in pairs(mergedTable) do
			if (state.personalMemories[i] == nil) or (state.personalMemories[i] ~= v) then
				send("rendCommandManager", "odinRendererDeleteMemoryMessage", SELF, i);
				if v.data.iconSkin ~= nil then
					send("rendCommandManager", "odinRendererNewMemoryMessage", SELF, i, v.eventName, v.data.icon, v.data.iconSkin, v.attributes.happiness, v.attributes.despair, v.attributes.fear, v.attributes.anger, v.attributes.importance, v.data.priority, false, parseDescription(v.data.description))
				else
					send("rendCommandManager", "odinRendererNewMemoryMessage", SELF, i, v.eventName, v.data.icon, "thoughtIcons.xml", v.attributes.happiness, v.attributes.despair, v.attributes.fear, v.attributes.anger, v.attributes.importance, v.data.priority, false, parseDescription(v.data.description))
				end
			end

			state.personalMemories[ i ] = v
			i = i+1
		end

		if i < 10 then
			for k=i+1,10 do
				send("rendCommandManager", "odinRendererDeleteMemoryMessage", SELF, k);
			end
		end
--]]

--[[
		local mSAStart = getTime()

		parseMemoriesForEmotions()

		VALUE_STORE["MSAUpdateTime"] =  VALUE_STORE["MSAUpdateTime"] + getTime() - mSAStart
		
		
		--]]

	end
	function createOriginStory()
		
		local storiesDB = EntityDB["HumanStats"].originStories
		local possibleStories = {}
		
		if SELF.tags["former_bandit"] then
			state.AI.strs["originStory"] = state.AI.name .. " left a life of Banditry to re-join civilization and become a productive member of society."
			return
		end
		
		if state.AI.strs["citizenClass"] == "Prisoner" then
			
			for k,story in pairs( storiesDB.byTrait["Prisoner"].any_class ) do
				possibleStories[#possibleStories + 1] = story
			end
			
			state.AI.traits["Prisoner"] = true
			state.AI.strs["originStory"] = state.AI.name .. " " .. possibleStories[ rand(1, #possibleStories) ]
			return
		end
		
		for trait,byclass in pairs(storiesDB.byTrait) do
			if state.AI.traits[trait] then
				for class,stories in pairs(byclass) do
					if class == "any_class" or
						class == state.AI.strs["socialClass"] then
						
						for k,story in pairs(stories) do
							possibleStories[#possibleStories + 1] = story
						end
					end
				end
			end
		end
		for skill,byskill in pairs(storiesDB.bySkill) do
			if state.AI.skills[skill] > 1 then
				for class,stories in pairs(byskill) do
					if class == "any_class" or
						class == state.AI.strs["socialClass"] then
						
						for k,story in pairs(stories) do
							possibleStories[#possibleStories + 1] = story
						end
					end
				end
			end
		end
		
		if #possibleStories < 3 then
			for class,stories in pairs(storiesDB.byClass) do
				if state.AI.strs["socialClass"] == class or
					class == "any_class" then
					
					for k,story in pairs(stories) do
						possibleStories[#possibleStories + 1] = story
					end
				end
			end
		end

		local diceRoll = rand(1, #possibleStories)
		state.AI.strs["originStory"] = state.AI.name .. " " .. possibleStories[diceRoll]
	end

	function makeMoodText()
		state.AI.strs["moodText"] = "Mood: ".. state.AI.strs.mood .. ". "
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", SELF.id, "mood_text", state.AI.strs.mood)
	end

 	function makeFriendsText()
 		-- add friends to descriptiveParagraph
          if #state.AI.friends > 0 then
               local friendString = "Friends: "
			 
               for i=1, #state.AI.friends do
				if state.AI.friends[i] and not state.AI.friends[i].deleted then
					
					local friendName = query(state.AI.friends[i],"getName")[1]
					
					if friendName then
						if i == #state.AI.friends or #state.AI.friends == 1 then
							-- last or only
							friendString = friendString .. friendName .. ". "
						else
							-- middle element
							friendString = friendString .. friendName .. ", "
						end
					end
				end
				state.AI.strs["friendsText"] = friendString
			end
          else
               state.AI.strs["friendsText"] = "Friends: sadly, none."  
          end
	end
	

	function makeSupplyText()
		-- add amount of supply to 1 paragraph
		--[[local tempSupplyString = ""
		local tempSupplyString2 = ""
		if state.AI.strs["socialClass"] == "lower" then
			local overseer = query ( "gameBlackboard", "gameObjectGetOverseerMessage", state.AI.currentWorkParty)[1]
			if overseer then
				local overseerAttribs = query(overseer, "getAIAttributes")[1]
				if overseerAttribs then
					local NCOSupplyCounter = overseerAttribs.ints["SupplyCounter"]
					if NCOSupplyCounter > 0 then
						tempSupplyString = "This unit's NCO has an active unit of supply."
					else
						tempSupplyString = "This unit's NCO does not have an active unit of supply."
					end
					if not overseerAttribs.bools.reserveSupply then
						tempSupplyString2 = "This unit's NCO has no reserve supplies. Make some at the armory for maximum combat effectiveness."
					else
						tempSupplyString2 = "This unit's NCO has some unused reserve supplies."
					end
				end
			end
			
		else
			if state.AI.ints["SupplyCounter"] > 0 then
				tempSupplyString = "This unit is currently benefiting from an active unit of supply!"
			else
				tempSupplyString = "This unit has no active units of supply."
			end
			if state.AI.bools["reserveSupply"]then
				tempSupplyString2 = "This unit has some reserve supplies and is ready for action."
			else
				tempSupplyString2 = "This unit's NCO has no reserve supplies. Make some at the armory for maximum combat effectiveness."
			end
		end
		
		state.AI.strs["supplyText"] =  tempSupplyString2 .. " " .. tempSupplyString]]
	end
	
	-- Set tags & workcrew filters for character class
	-- This is a subset of "changeCharacterClass" that can be called specifically
	-- to reset tags/filters without changing profession.
	function setOverseerWorkcrewFilters()

		local overseer = query("gameBlackboard",
						   "gameObjectGetOverseerMessage",
						   state.AI.currentWorkParty)[1]
				
		if state.AI.strs["socialClass"] == "middle" then
			
			if state.AI.strs["citizenClass"] == "NCO" or
				state.AI.strs["citizenClass"] == "Militia NCO" then

				send("gameBlackboard", "gameSetWorkPartyMilitary", overseer, true)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "hauling", false)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "construction", false)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "farming", false)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "foraging", false)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "mining", false)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "chopping", false)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "workshop", false)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "hunting", false)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "military", true)
	
				state.AI.bools["reserveSupply"] = true
				state.AI.ints["SupplyCounter"] = 0
				SELF.tags["military"] = true
				makeSupplyText()
				
			elseif state.AI.claimedWorkBuilding then

				send("gameBlackboard", "gameSetWorkPartyMilitary", overseer, false)
				
				local building_tags = query(state.AI.claimedWorkBuilding,"getTags")[1]
				if building_tags.workshop then 
					send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "hauling", true)
				else
					send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "hauling", false)
				end
				
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "construction", false)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "farming", false)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "foraging", false)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "mining", false)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "chopping", false)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "workshop", true)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "hunting", false)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "military", false)
				
			else
				-- if overseer, then filters to whatever.
				-- UPDATE: reset filters.
				
				send("gameBlackboard", "gameSetWorkPartyMilitary", overseer, false)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "hauling", true)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "construction", true)
				--send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "farming", true)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "foraging", true)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "mining", true)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "chopping", true)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "workshop", true)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "hunting", false)
				send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "military", false)
			end
		end
	end
	
	function setProfessionTags()
		
		-- first, remove all possible job tags.
		for k,v in pairs(EntitiesByType["workshop"]) do
			if v.job_tags then
				for a,b in pairs(v.job_tags) do
					SELF.tags[b] = nil	
				end
			end
			
			if state.AI.strs["socialClass"] == "lower" then
				if v.job_tags_lc then
					for a,b in pairs(v.job_tags_lc) do
						SELF.tags[b] = nil	
					end
				end
			elseif state.AI.strs["socialClass"] == "middle" then
				if v.job_tags_mc then
					for a,b in pairs(v.job_tags_mc) do
						SELF.tags[b] = nil	
					end
				end
			end
		end
		
		for k,v in pairs(EntitiesByType["office"]) do
			if v.job_tags then
				for a,b in pairs(v.job_tags) do
					SELF.tags[b] = nil	
				end
			end
			
			if state.AI.strs["socialClass"] == "lower" then
				if v.job_tags_lc then
					for a,b in pairs(v.job_tags_lc) do
						SELF.tags[b] = nil	
					end
				end
			elseif state.AI.strs["socialClass"] == "middle" then
				if v.job_tags_mc then
					for a,b in pairs(v.job_tags_mc) do
						SELF.tags[b] = nil	
					end
				end
			end	
		end
		
		state.AI.ints.sermon_timer = nil

		if state.AI.strs["socialClass"] == "lower" then
			local overseer = false
			if state.AI.currentWorkParty then 
				overseer = query("gameBlackboard","gameObjectGetOverseerMessage",state.AI.currentWorkParty)[1]
			end
			
			if overseer then
				local claimedBuilding = query(overseer,"getClaimedWorkBuilding")[1]
				if claimedBuilding then
					local building_name = query(claimedBuilding,"getBuildingName")[1]
					local building_data = EntityDB[building_name]
					if building_data.job_tags then
						for k,v in pairs(building_data.job_tags) do
							SELF.tags[v] = true	
						end
					end
					
					if building_data.job_tags_lc then
						for k,v in pairs(building_data.job_tags_lc) do
							SELF.tags[v] = true	
						end
					end
				end
				
				local overseer_tags = query(overseer, "getTags")[1]
				
				if overseer_tags and overseer_tags["military"] then
					if state.AI.skills.militarySkill < 2 then
						SELF.tags.militia = true
					end
				end
			end
			
		elseif state.AI.strs.socialClass == "middle" then
			
			local workBuildingTags = false
			if state.AI.claimedWorkBuilding then
				workBuildingTags = query(state.AI.claimedWorkBuilding,"getTags")[1]
				SELF.tags.workshop_jobs = true
				
				local building_name = query(state.AI.claimedWorkBuilding,"getBuildingName")[1]
				local building_data = EntityDB[building_name]
				if building_data.job_tags then
					for k,v in pairs(building_data.job_tags) do
						SELF.tags[v] = true	
					end
				end
				if building_data.job_tags_mc then
					for k,v in pairs(building_data.job_tags_mc) do
						SELF.tags[v] = true	
					end
				end
			end
				
			if state.AI.strs["citizenClass"] == "NCO" then
				SELF.tags.military = true
				SELF.tags.military_decisiontree = true

				-- if assigned to barracks, do military training.
				if workBuildingTags and workBuildingTags.barracks then
					
					SELF.tags.barracks_jobs = true
				end
			elseif state.AI.strs["citizenClass"] == "Militia NCO" then
				-- if assigned to barracks, do military training.
				if workBuildingTags and workBuildingTags.barracks then
					SELF.tags.barracks_jobs = true
				end
				SELF.tags.militia = true
				
			elseif state.AI.strs["citizenClass"] == "Vicar" then
				state.AI.ints.sermon_timer = 0
			end
		end
	end
	
	-- give us a new character class
	function changeCharacterClass( newclass )
		if state.lockCharacterClass == true or SELF.tags.dead then
			return
		end
		
		-- STEP 1:
		-- if no class given, select correct class.
		local overseer = false
		
		if state.AI.strs["socialClass"] == "lower" and state.AI.currentWorkParty then 
			overseer = query( "gameBlackboard",
						    "gameObjectGetOverseerMessage",
						    state.AI.currentWorkParty)[1]
		end
          
		-- Tutorial Block.
          if query("gameSession", "getSessionBool", "caseTutorialActive")[1] == true then 

               if state.AI.strs["socialClass"] == "lower" then
                    if overseer then
                         local director = query("gameSession","getSessiongOH","event_director_tutorial")[1]
                         send(director, "setKeyString", "workerAssigned", "yes")
                         
                         local overseer_tags = query(overseer, "getTags")[1]
                         if overseer_tags and (overseer_tags["military"] or overseer_tags.militia) then
                              local safety = query(overseer, "getClaimedWorkBuilding") -- this is to make sure you're not using the starting NCO.
                              if (safety) and (safety[1]) then
                                   send(director, "setKeyString", "redcoatTrained", "yes")
                              end
                         end
                    end
               end
          end
		
		local building_name = nil
		local building_data = nil
		local building = nil
		if overseer then
			building = query(overseer,"getClaimedWorkBuilding")[1]
			if building then
				building_name = query(building,"getName")[1]
				building_data = EntityDB[building_name]
			end
		elseif state.AI.strs["socialClass"] == "middle" then
			if state.AI.claimedWorkBuilding then
				building_name = query(state.AI.claimedWorkBuilding,"getName")[1]
				building_data = EntityDB[building_name]
			end
		end
		
		if newclass == "" then
			if state.AI.strs["socialClass"] == "lower" then
				newclass = "Labourer"
				if overseer then
					
					if building_data then
						if building_data.labourer_class then
							newclass = building_data.labourer_class
						end
					end
					
					local overseer_tags = query(overseer, "getTags")[1]
					if overseer_tags and overseer_tags.military then
						if query("gameSession", "getSessionBool", "military_training4_unlocked")[1] then
							newclass = "Footsoldier"
						else
							newclass = "Militia Footsoldier"
						end
					end

					if not newclass then
						-- assign character class based on highest skill of overseer.
						local result = query(overseer, "getAIAttributes")[1]
						local overseerHighestStat = -1
						local overseerHighestSkillName = " "
						for k,v in pairs(result.skills) do
							if v > overseerHighestStat then
								overseerHighestSkillName = k
								overseerHighestStat = v
							end
						end
					end
					
				--else
					-- no overseer, presumed dead = stay soldier
					--[[if state.AI.strs["citizenClass"] == "Footsoldier" then
						newclass = "Footsoldier"
					elseif state.AI.strs["citizenClass"] == "Militia Footsoldier" then
						newclass = "Militia Footsoldier"
					end]]
				end
				
			elseif state.AI.strs["socialClass"] == "middle" then
				-- claimed a work building? then find class from that.
				
				if state.AI.claimedWorkBuilding then

					send("rendCommandManager", "uiSetWorkPartyType", SELF, 1)
					
					if building_data and building_data.overseer_class then
						newclass = building_data.overseer_class
					else
						newclass = "Artisan"
					end
					
					if building_name == "Barracks" then
						send("rendCommandManager", "uiSetWorkPartyType", SELF, 2)
						if state.AI.skills.militarySkill > 1 then
							newclass = "NCO"
						else
							newclass = "Militia NCO"
						end
					end
					
				else
					newclass = "Overseer"
					send("rendCommandManager", "uiSetWorkPartyType", SELF, 0)
				end
				-- should class change based on character skill?
				-- perhaps not ... yet?
			end
		end
		
		if newclass == "Militia NCO" and
			state.AI.skills.militarySkill > 1 then
				
			newclass = "NCO"
		elseif newclass == "NCO" and
			state.AI.skills.militarySkill < 2 then
			
			newclass = "Militia NCO"
		end
		
		printl("ai_agent", state.AI.name .. " is switching class from: " ..
			  state.AI.strs["citizenClass"] .. " to: " .. tostring(newclass) )
		
		-- soldiers  lose all training points when assigned off-duty.
		if state.AI.strs["citizenClass"] == "Footsoldier" or
			(state.AI.strs["citizenClass"] == "Militia Footsoldier" and
			newclass ~= "Footsoldier") then
			
			state.AI.ints.militaryTrainingPoints = 0
		end
		
		if newclass ~= "Labourer" and
			state.AI.strs["citizenClass"] == newclass then
		end

		-- STEP 4
		-- setting some tags; but the more fundamental ones.
		SELF.tags.industrial_artisan = nil
		SELF.tags[ state.AI.strs["citizenClass"] ] = nil
		SELF.tags[ string.lower( state.AI.strs["citizenClass"]) ] = nil
		SELF.tags[ state.AI.strs["socialClass"] ] = nil
		SELF.tags[ state.AI.strs["socialClass"] .. "_class" ] = nil

		state.AI.strs["citizenClass"] = newclass
		
		local entityData = EntityDB[ state.AI.strs["citizenClass"] ]		
		state.AI.strs["socialClass"] = entityData.socialClass
		
		SELF.tags[ state.AI.strs["citizenClass"] ] = true
		SELF.tags[ string.lower( state.AI.strs["citizenClass"]) ] = true
		SELF.tags[ state.AI.strs["socialClass"] ] = true
		SELF.tags[ state.AI.strs["socialClass"] .. "_class" ] = true
		
		if entityData.military == 1 then
			state.AI.strs["loadout_tool"] = "firearm"
			--Set up new HP max if tech agrees
			local maxHealthBonus = query("gameSession", "getSessionInt", "militaryHealthBonus")[1]
			if maxHealthBonus ~= 0 then
				local humanstats = EntityDB["HumanStats"]
				state.AI.ints["healthMax"] = humanstats["healthMax"] + maxHealthBonus
			end
		else
			state.AI.strs["loadout_tool"] = ""
			
			--Set up new HP max if tech agrees
			local maxHealthBonus = query("gameSession", "getSessionInt", "civilianHealthBonus")[1]
			if maxHealthBonus ~= 0 then
				local humanstats = EntityDB["HumanStats"]
				state.AI.ints["healthMax"] = humanstats["healthMax"] + maxHealthBonus
			end
		end

		-- Set up model.
		local models = getModelsForClass( state.AI.strs["citizenClass"],
								   state.AI.strs["gender"],
								   state.AI.strs["variant"] )
		state.models = models
		state.animSet = models["animationSet"]
		
		local currentAnim = "idle"
		if SELF.tags.sleeping then

			if state.AI.strs["sleepLocation"] == "ground" or
				state.AI.strs["sleepLocation"] == "building" then
				currentAnim = "sleep_loop"
			elseif state.AI.strs["currentBedType"] and
				state.AI.strs["currentBedType"] ~= "" then
				
				currentAnim = "sleep_" .. state.AI.strs["currentBedType"] .. "_" .. state.AI.strs["currentBedDirection"] .. "_loop"
			end
		end
		
		if not SELF.tags.in_cult_robes then
			local hatmodel = ""

			send("rendOdinCharacterClassHandler",
				"odinRendererSetCharacterGeometry", 
				state.renderHandle,
				models["torsoModel"], 
				"",
				hatmodel,
				models["animationSet"],
				currentAnim)
		end
		
		-- UI setup.
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle,
			"occupation",
			state.AI.strs["citizenClass"])
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle,
			"socialClass",
			state.AI.strs["socialClass"])
		
		if state.AI.strs["socialClass"] == "upper" then
			send("rendUIManager", "uiRemoveColonist", SELF.id)
			-- send("rendUIManager", "uiRemoveWorker", SELF.id) -- causes scripterror.
			send("gameBlackboard", "gameObjectRequestNewOverseerMessage", nil, 0)
			SELF.tags["workshop_jobs"] = nil
			SELF.tags["work_party_jobs"] = nil
			send(SELF, "gameObjectSetOverseerMessage",
				SELF,
				nil,
				nil )
			--state.AI.currentWorkParty = nil
		end

		setProfessionTags()
		setOverseerWorkcrewFilters()
		
		send(SELF,"updateWorkQoL") -- for new office or lack of office.
		
		-- oh hey, setup weapons.
		
		local profession = state.AI.strs["citizenClass"]
		local data = EntityDB[profession]
		if data.ranged_weapons then
			-- set random? default weapon
			send(SELF,"setWeapon","ranged",data.ranged_weapons[ rand(1,#data.ranged_weapons)] )
		end
		
		if data.melee_weapons then
			local random_weapon = data.melee_weapons[ rand(1,#data.melee_weapons)]
			send(SELF,"setWeapon","melee",random_weapon)
		else
			send(SELF,"setWeapon","melee","default") 
		end
		
		if SELF.tags.military then
			-- do ranged weapon from barracks.
			if state.AI.strs["socialClass"] == "lower" and state.AI.currentWorkParty then 
				local overseer = query("gameBlackboard","gameObjectGetOverseerMessage",state.AI.currentWorkParty)[1]
				local rax = query(overseer, "getClaimedWorkBuilding")[1]
				if rax then
					-- get loadout weapon from barracks
					local weapon = query(rax,"getWeaponLoadout")[1]
					if state.AI.strs.ranged_weapon ~= weapon then
						send(SELF,"setWeapon","ranged",weapon)
					end
				end
			elseif state.AI.strs["socialClass"] == "middle" then
				if state.AI.claimedWorkBuilding then
					-- get loadout weapon from barracks
					local weapon = query(state.AI.claimedWorkBuilding,"getWeaponLoadout")[1]
					if state.AI.strs.ranged_weapon ~= weapon then
						send(SELF,"setWeapon","ranged",weapon)
					end
				end 
			end
		end
		
	end

	function setposition( x, y )
		newPos = gameGridPosition:new()
		newPos:set( x, y )
		state.AI.position = newPos
	end

	function setReligionPhrase()
		state.AI.strs["religionPhrase"] = "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN worships at the altar of celestial order."
	end

	function setPoliticalPrefPhrase()
		state.AI.strs["politicalPrefPhrase"] = "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN has pledged allegiance to the Queen."
	end

	function setAgePhrase()
		local diceroll = rand(15,45)
		local ageOfTheQueen = 70
		state.AI.ints["age"] = diceroll
		local year = ageOfTheQueen - diceroll
		local numberEnding = "th"
		if year%10 == 1 then
			numberEnding = "st"
		elseif year%10 == 2 then
			numberEnding = "nd"
		elseif year%10 == 3 then
			numberEnding = "rd"
		end
		state.AI.strs["agePhrase"] = "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN was born in the " .. year .. numberEnding .. " year of the Reign of the Queen."
	end

	function setPhysicalTraitsPhrase()
		-- TODO set these based on actual character traits
		-- For now, just having some fun here.
		local physicalTraits = { " is reasonably hardy though looks suspiciously like the other settlers.",
										" is a sturdy subject and has one of those faces you just keep seeing everywhere.",
										" has a good colonial build though walks with a suspicious gait.",
										" is a fine specimen with a stance neither above nor below that which is appropriate.",
										" has a hungry look in the eyes though appears otherwise capable." }	
		state.AI.strs["physicalTraitsPhrase"] = "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN" .. physicalTraits[ rand(1,5) ]
	end

	function recalcShiftLength()

	-- need to update happiness dependent shifts with recalcShiftLength, but nothing else
	-- currently (2016-5-27) we have 2 happiness dependent shifts, first one enabled is 2, disabled is 3 & threshold is 5 shifts
	-- second one is enabled on 4, disabled on 5, threshold is 6 shifts
	
	-- WARNING 2016-7-15: This code is possibly depreciated, shift stuff is apparently handled automagically on the C++ side. I wouldn't touch it. -ct
		if SELF.tags.middle_class then
			if state.AI.ints.happiness then 
				for k,v in ipairs(EntityDB.HumanStats.happinessToWorkshifts) do
					
					if state.AI.ints.happiness > v.happiness then

						for i=1,8 do
							tempVar = state.AI.ints["hour" .. i]
							if tempVar == 2 or tempVar == 3 then
								if v.shifts > 4 then
									state.AI.ints["hour" .. i] = 2
								else
									state.AI.ints["hour" .. i] = 3
								end
							elseif tempVar == 4 or tempVar == 5 then
								if v.shifts > 5 then
									state.AI.ints["hour" .. i] = 4
								else
									state.AI.ints["hour" .. i] = 5
								end
							end
						end
						break
					end
				end
			end
			send("rendCommandManager", 
					"gameSetWorkPartyWorkShift", 
					SELF, 
					state.AI.ints.hour1, 
					state.AI.ints.hour2, 
					state.AI.ints.hour3, 
					state.AI.ints.hour4, 
					state.AI.ints.hour5, 
					state.AI.ints.hour6, 
					state.AI.ints.hour7, 
					state.AI.ints.hour8)
		end
	end

	function setDescriptiveParagraph()

		local skillDescriptionStrings = {}

		state.descriptiveParagraph = 		""
		
		state.descriptiveParagraph = state.descriptiveParagraph .. state.AI.strs.dailyJournalText
		state.descriptiveParagraph = state.descriptiveParagraph .. "\n"
		state.descriptiveParagraph = state.descriptiveParagraph .. "\n" .. state.AI.strs.weeklyJournalText
		state.descriptiveParagraph = state.descriptiveParagraph .. "\n"
					
		state.descriptiveParagraph = state.descriptiveParagraph .. "\n" .. state.AI.strs["friendsText"] .. " \n"
		state.descriptiveParagraph = state.descriptiveParagraph .. "\n"
			.. state.AI.strs["originStory"] .. " "
			--.. state.AI.strs["physicalTraitsPhrase"] .. " \n"
			-- .. state.AI.strs["religionPhrase"] .. " \n"
			-- .. state.AI.strs["politicalPrefPhrase"] .. " \n"
			.. state.AI.strs["agePhrase"] .. " \n"
			

		local hungerString = ""
		if state.AI.ints.hunger == 0 then
			hungerString = "Not hungry."
		elseif state.AI.ints.hunger < 2 then
			hungerString = "A bit hungry."
		elseif state.AI.ints.hunger < 4 then
			hungerString = "Rather hungry."
		elseif state.AI.ints.hunger < 6 then
			hungerString = "Extremely hungry."
		elseif state.AI.ints.hunger < state.AI.ints["starvationDeathHunger"] then
			hungerString = "Starving to death."
		end

		local tirednessString = ""
		if state.AI.ints.tiredness == 0 then
			tirednessString = "Awake and alert."
		elseif state.AI.ints.tiredness < 2 then
			tirednessString = "Tired and looking for a bed."
		elseif state.AI.ints.tiredness < 3 then
			tirednessString = "Tired enough to sleep anywhere warm."
		else
			tirednessString = "So tired they're willing to sleep anywhere."
		end

		--[[
		state.descriptiveParagraph = state.descriptiveParagraph .. "\n" ..
			"Hunger: " .. hungerString

		state.descriptiveParagraph = state.descriptiveParagraph .. "\n" ..
			"Tiredness: " .. tirednessString ..
			"\n" .. "\n"


			for i=1, #skillDescriptionStrings do
				state.descriptiveParagraph = state.descriptiveParagraph .. skillDescriptionStrings[i] .. " \n"
			end
		--]]
	end

	>>

	state
	<<
		gameAIAttributes AI
		string vocalID
		string animSet
		int renderHandle
		string descriptiveParagraph
		string seat
		bool release
		bool lockCharacterClass
		string changeToCharacterClass
	>>

	receive Create( stringstringMapHandle init )
	<<
		local entityName = init["legacyString"]
		printl("ai_agent", "placing: " .. entityName)
		
          local isMurderer = false
		local isCriminal = false
		local isFormerBandit = false
		
          if entityName == "murderer" then
               isMurderer = true
               entityName = "citizen"
          elseif entityName == "Prisoner" then
			isCriminal = true
			state.lockCharacterClass = true
		elseif entityName == "Naturalist" then
			-- Naturalist granted via special event
			-- keep this job class no matter what!
			state.lockCharacterClass = true
		elseif entityName == "Bandit" then
			isFormerBandit = true
			state.AI.traits["Former Bandit"] = true
		elseif entityName == "lower_class" then
			entityName = "Labourer"
		elseif entityName == "middle_class" then
			entityName = "Overseer"
		end
		
		if init.force_pioneer then
			state.AI.traits["Pioneering Spirit"] = true
		end
		
		if init.lockCharacterClass then
			if init.lockCharacterClass == "true" then
				state.lockCharacterClass = true
			end	
		end
		
		if entityName == "citizen" then
			--  make it totally random
			local variousClasses = {
				--"Naturalist",
				"Overseer",
				"Labourer",
				"Labourer",
				"Labourer",
				"Overseer",
				"Labourer",
				"Labourer",
				"Labourer",
				"Labourer",
				--"NCO",
				--"Footsoldier",
				"Poet",
				"Aristocrat",
				"Capitalist",
				--"Vicar", }
				}
			
			entityName = variousClasses[ rand(1, #variousClasses) ]
		end
		local maxHealthBonus = query("gameSession", "getSessionInt", "civilianHealthBonus")[1]
		if entityName == "Footsoldier" or entityName == "NCO" then
			send("gameSession", "incSessionInt", "militaryCount", 1)
			maxHealthBonus = query("gameSession", "getSessionInt", "militaryHealthBonus")[1]
		end
		
		state.release = false
		
		state.AI.bools.first_placement = true -- used by gameobjectplace.
		
		local humanstats = EntityDB["HumanStats"]
		local worldstats = EntityDB["WorldStats"]

		state.AI.ints["3s_per_inebriation"] = divCeil(divCeil(worldstats["dayNightCycleSeconds"]*3, 3),10)  -- *3 to greatly reduce rate of booze consumption; was:  /3 / 10  )

		state.AI.ints["inebriation_counter"] = rand( 1, state.AI.ints["3s_per_inebriation"] )
		state.AI.ints["idealSleepTime"] = div( worldstats.nightTimeSeconds, 3.2 ) -- reduced further. -- 80% of the night -- (I changed it because people sleep so long -DJ)
		
		state.AI.currentParty = nil
		
		state.AI.ints.militaryTrainingPoints = 0
		state.AI.ints.militaryTrainingPointsGoal = humanstats.militaryTrainingPointsToRedcoat
		
		state.AI.ints["jobPriorityVal"] = 6 -- used for the combat decision tree.
		state.AI.bools["is_married"] = false
		state.AI.ints["corpse_timer"] = humanstats["corpseRotTimeDays"] * worldstats["dayNightCycleSeconds"] * 10 -- in gameticks
		state.AI.ints["corpse_vermin_spawn_time_start"] = div(state.AI.ints["corpse_timer"],2)
		state.AI.ints["fire_timer"] = 50 -- time between sending out "set on fire" pulses
		state.AI.ints["alarmTimer"] = 11
		state.AI.ints["emotionAnimationTimer"] = 0
		state.AI.ints["starvationDeathHunger"] = humanstats["starvationDeathDays"]
		state.AI.ints["hunger"] = 0 -- rand(0,7) -- starting hunger
		state.AI.ints["tiredness"] = 0 --rand(0,3)  -- starting tiredness
		
		state.AI.ints.lastDayInCombat = -1 -- used for safety qol
		
		if isCriminal or isFormerBandit then
			-- they aren't treated well
			state.AI.ints["hunger"] = rand(0,4)
		end
		
		if init.hunger then
			state.AI.ints.hunger = tonumber( init.hunger )
		end
		
		if init.tiredness then
			state.AI.ints.tiredness = tonumber( init.tiredness )
		end
		state.lastSleepDay = -1

		state.AI.ints["despair"] = 1
		state.AI.ints["strongestEmotionValue"] = 0 
		--state.AI.ints["morale"] = 0
		state.AI.ints["deathsWitnessed"] = 0
		state.AI.strs["mood"] = "indifferent"
		state.AI.strs["oldMood"] = ""
		state.AI.strs["lastThought"] = "indifferent"
		state.AI.strs["friendsText"] = "Friends: Sadly, none."
		state.AI.strs["moodText"] = "Huh, weird, this character has no mood at all!"
		--state.AI.strs["moraleText"] = "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN will stand firm against enemies."
		state.AI.strs["despairText"] = "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN is of stable mind."
		state.AI.strs["ammoText"] = " "
		state.AI.strs["supplyText"] = " "
		state.AI.strs["dailyJournalText"] = "Yesterday: No journal entries yet."
		state.AI.strs["weeklyJournalText"] = "Last Week: No journal entries yet."

		state.AI.ints.hoursPerShift = 0 --populating these so requirements blocks don't freak out if we don't have one right away.
		state.AI.ints.hoursPerShiftLocked = -1

		state.AI.ints.hour1 = 0 -- 0: enabled, 1: disabled, 2: first happiness-enabled hour ENABLED (lowest threshold) , 3: first happiness-enabled hour DISABLED (lowest threshold), 4: second happiness-enabled hour ENABLED, 5: second happiness-enabled hour DISABLED
		state.AI.ints.hour2 = 0
		state.AI.ints.hour3 = 2
		state.AI.ints.hour4 = 1
		state.AI.ints.hour5 = 0
		state.AI.ints.hour6 = 0
		state.AI.ints.hour7 = 4
		state.AI.ints.hour8 = 1

		state.AI.ints["inebriation"] = 0 -- physiological effects of alcohol, 1-10 (per day) scale
		state.AI.ints["SupplyCounter"] = 0

		state.AI.ints["health"] = humanstats["healthMax"] -- yeah, start at max
		state.AI.ints["healthMax"] = humanstats["healthMax"] + maxHealthBonus -- adjusts based on tech stuff
		state.AI.ints["healthRegenTimer"] = humanstats["healthRegenTimerSeconds"]-- in seconds
		
		if init.numAfflictions then
			state.AI.ints["numAfflictions"] = tonumber( init.numAfflictions )
		else
			state.AI.ints["numAfflictions"] = 0
        end

		state.AI.ints["fire_timer"] = 10
		state.AI.ints["grenades"] = 0
		state.AI.ints["grenadesMax"] = humanstats["grenadesMax"]

		state.AI.ints["inCombatTimer"] = 0
		
		state.AI.stamina = 100

		state.AI.bools["isRottedCorpse"] = false
		state.AI.bools["dilligent"] = true

		state.AI.ints.revoltTimer = 0

		state.AI.ints["partner"] = 0
		state.AI.ints["rival"] = 0
		state.AI.ints["maxNumberFriends"] = humanstats["maxNumberFriends"]
		state.AI.ints["relationshipThresholdFriend"] = humanstats["relationshipThresholdFriend"] -- double this after each friend is made
		state.AI.ints["relationshipThresholdMarriage"] = humanstats["relationshipThresholdMarriage"]

		for i=1, #EntityDB.HumanStats.skillNameList do
			state.AI.skillDisplayNames[i] = humanstats.skillDisplayNameList[i]
		end

		for i=1,#EntityDB.HumanStats.skillNameList do
			if init["skill_" .. humanstats.skillNameList[i] ] then
				state.AI.skills[ humanstats.skillNameList[i]] = tonumber( init["skill_" .. humanstats.skillNameList[i] ] ) 
			else
				state.AI.skills[ humanstats.skillNameList[i] ] = 1
			end
			state.AI.skillEvents[ humanstats.skillNameList[i]] = 0
			state.AI.skillLevelDisplayText[ humanstats.skillNameList[i]] = humanstats.skillLevelStrings[ state.AI.skills[ humanstats.skillNameList[i]] ]
		end
		
		state.AI.bools[ "isWieldingTool" ] = false
		
		if init.gender then
			state.AI.strs.gender = init.gender
			state.AI.strs.genderstr = init.gender
		else
			if( rand(0,100) < 50 ) then 
				state.AI.strs["gender"] = "male"
				state.AI.strs["genderstr"] = "male"
			else
				state.AI.strs["gender"] = "female"
				state.AI.strs["genderstr"] = "female"
			end
		end
		
		if state.AI.strs.gender == "male" then
			state.AI.strs.CAPITALIZED_SUBJECTIVE = "He"
			state.AI.strs.SUBJECTIVE = "he"
			state.AI.strs.CAPITALIZED_POSSESSIVE = "His"
			state.AI.strs.POSSESSIVE = "his"
			state.AI.strs.CAPITALIZED_OBJECTIVE = "Him"
			state.AI.strs.OBJECTIVE = "him"
			
		elseif state.AI.strs.gender == "female" then
			state.AI.strs.CAPITALIZED_SUBJECTIVE = "She"
			state.AI.strs.SUBJECTIVE = "she"
			state.AI.strs.CAPITALIZED_POSSESSIVE = "Her"
			state.AI.strs.POSSESSIVE = "her"
			state.AI.strs.CAPITALIZED_OBJECTIVE = "Her"
			state.AI.strs.OBJECTIVE = "her"
		end
		--[[
				state.AI.strs.CAPITALIZED_SUBJECTIVE = "They"
				state.AI.strs.SUBJECTIVE = "they"
				state.AI.strs.CAPITALIZED_POSSESSIVE = "Their"
				state.AI.strs.POSSESSIVE = "their"
				state.AI.strs.CAPITALIZED_OBJECTIVE = "Them"
				state.AI.strs.OBJECTIVE = "them"
				
			]]
			
		state.AI.ints["pref"] = 0 -- 1=m, 2=f, 3=m/f, 0=none
		local r = rand(1,100)
		if state.AI.strs["gender"] == "female" then
			if r < 11 then
				state.AI.ints["pref"] = 2
			elseif r < 12 then
				state.AI.ints["pref"] = 3
			elseif r < 100 then
				state.AI.ints["pref"] = 1
			end
		else
			if r < 11 then
				state.AI.ints["pref"] = 1
			elseif r < 12 then
				state.AI.ints["pref"] = 3
			elseif r < 100 then
				state.AI.ints["pref"] = 2
			end
		end
		
		if init.variant then state.AI.strs["variant"] = init.variant end
	
		function getFirstName()
			if( state.AI.strs["genderstr"] == "male" ) then
				return maleFirstNames[ rand(1, #maleFirstNames ) ]
			else
				return femaleFirstNames[ rand(1, #femaleFirstNames ) ]
			end
		end
		
		function getLastName()
			local lastname = ""
			
			local lastA = lastnamePrefixes[ rand(1, #lastnamePrefixes ) ]
			local lastB = lastnameRoots[ rand(1, #lastnameRoots ) ]
			local lastC = lastnameSuffixes[ rand(1, #lastnameSuffixes ) ]
	
			local temprand = rand(1,9)
			if( temprand == 1 ) then
				lastname = lastA .. lastB
			elseif (temprand == 2) then
				lastname = lastB .. lastC
			elseif (temprand == 3) then
				lastname = lastB
			elseif (temprand == 4) then
				--state.AI.strs["lastName"] = lastnames[ rand(1, #lastnames ) ]
				lastname = lastA .. lastC
			else
				--state.AI.strs["lastName"] = lastA .. lastB .. lastC
				lastname = lastnames[ rand(1, #lastnames ) ]
			end
	
			lastname = string.upper( lastname:sub(1,1) ) .. lastname:sub(2, #lastname)
			return lastname
		end
			
		if init.firstName then
			state.AI.strs["firstName"] = init.firstName
		else
			state.AI.strs["firstName"] = getFirstName()
		end

		-- TODO: (social) class-based naming
		if init.lastName then
			state.AI.strs["lastName"] = init.lastName
		else
			state.AI.strs.lastName = getLastName()
			while string.len(state.AI.strs.firstName .. " " .. state.AI.strs.lastName) > humanstats.nameLengthMax do
				state.AI.strs["lastName"] = getLastName()
			end
		end
		
		state.AI.name = state.AI.strs["firstName"] .. " " .. state.AI.strs["lastName"] 

		if init.class then
			state.AI.strs["citizenClass"] = init.class
		else
			if string.len(entityName) ~= 0 then	
				state.AI.strs["citizenClass"] = entityName
			else
				state.AI.strs["citizenClass"] = playableClasses[ rand(1, #playableClasses ) ]
			end
		end
		
		entityData = EntityDB[ state.AI.strs["citizenClass"] ]
		
		setposition( 0, 0 )
		
		local models = getModelsForClass( state.AI.strs["citizenClass"],
								   state.AI.strs["gender"],
								   "" )
		
		-- hope this works
		if not state.AI.strs["variant"] then
			state.AI.strs["variant"] = models["variant"]
		else
			models["variant"] = state.AI.strs["variant"]
		end
		
		if init.headModel then models.headModel = init.headModel end
		if init.hairModel then models.hairModel = init.hairModel end
			
		local hatmodel = models["hatModel"]
		local hairmodel = models["hairModel"]
		if state.AI.strs["citizenClass"] == "Prisoner" then
			hatmodel = "models/hats/prisonerHat.upm"
			hairmodel = ""
		--[[else
			-- TEST.
			hatmodel = "models/hats/selenianHeadBlob.upm"
			if rand(1,2) == 1 then
				hatmodel = "models/hats/selenianHeadStalk.upm"
			end]]
		end
		
		state.models = models

		-- CECOMMPATCH - HATS!
			local hatsel = {
				"hatCapotainHatBlack", --1 PILGRIM - ALL
				"hatCapotainHatBrown", --2 PILGRIM - ALL
				
				"vicarHat", --3 l
				"strawSunHatPointed", --4 l
				"strawSunHatRound", --5 l
				"hatSailor02", --6 l
				
				"boater", --7 l/m
				"bossOfThePlains", --8 l/m
				"bowler", --9 l/m
				"deerstalker", --10 l/m
				"flatCap", --11 l/m
				
				"hatSailor01", --12 m
				"hatBergstromHat", --13 m
				"hatAlpineGreen00", --14 m
				"hatAlpineGreen01", --15 m
				"pithHelmet", --16 m
				"fez", --17 m 
				
				"tophat", --18 u
				"merchantHat", --19 u
				"peakedcap" --20 u
				
				--"phrygiancap", -- silly?
				--"ushankaBrown", -- silly?
				--"ushankaGrey", -- silly?
				
				--"wushaCap" -- silly
				--"prisonerHat", -- silly
				--"chefHat", -- silly
				--"bearSkinHat", -- silly
				
				--"bicornhat", -- military
				--"tricornBrown", -- military
				--"tricornNavy", -- military
				--"stahlmarkianShako", -- military
				--"morionHelmet00", -- military
				--"occultInspectorHC", -- military
				--"occultInspectorLC", -- military
				--"kepi", -- military
				--"kepiSlouch", -- military
				--"shako", -- military
				--"pickelhaubel", -- military
				
				--"pufferHelm", -- fishpeople
				--"selenianHeadBlob", -- fishpeople
				--"selenianHeadStalk", -- firshpeople
				--"whelkGalea", -- fishpeople
				
				--"gogglesWornOnFace", -- floats above head
			}
			
			-- male-only because 1) to prevent a lot of bald females, and 2) because females already have hats attached to head models.
			if state.AI.strs["gender"] == "male" then
				-- chance to have a hat at all
				if rand(1,100) <= 40 then
					-- artistocrats wouldn't be caught dead in headwear below their class
					if entityName == "Aristocrat" then
						hatchoice = rand(18,20)
					else
						-- should we split by social class, or use the pilgrim hat?
						if rand(1,100) <= 10 then
							if entityName == "Overseer" then
								hatchoice = rand(7,17)
							else
								hatchoice = rand(3,11)
							end
						else
							-- default to pilgrim hat, just because it's awesome
							hatchoice = rand(1,2) -- 2 versions of it
						end
					end
					
					-- create the model string from our convoluted randomization
					hatmodel = "models/hats/" .. hatsel[hatchoice] .. ".upm"
					hairmodel = "" -- gotta hide the hair to prevent clipping
					printl("CECOMMPATCH - HAT: " .. hatsel[hatchoice])
					
					-- TO DO: manually review which hair/hat combos actually clip to prevent tons of bald folks, as well as give females a chance to be habbidashered
				end
			end
		-- /CECOMMPATCH

		send("rendOdinCharacterClassHandler",
			"odinRendererCreateCitizen", 
			SELF, 
			models["torsoModel"], 
			models["headModel"],
			hairmodel, 
			hatmodel, 
			models["animationSet"], 0, 0 )
		
		state.animSet = models["animationSet"]
		
		send("rendOdinCharacterClassHandler",
			"odinRendererCharacterSetBooleanAttribute",
			state.renderHandle,
			"automatic_interaction",
			true)

		send("rendOdinCharacterClassHandler",
			"odinRendererFaceCharacter", 
				state.renderHandle, 
				state.AI.position.orientationX,
				state.AI.position.orientationY )

		-- To differentiate these from non player controlled characters
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage",
					state.renderHandle, "characterClass", "citizen")


		if entityName == "Naturalist" then
			-- they get a boost to sight radius
			state.AI.ints["sightRadius"] = humanstats["sightRadius"] * 1.5
		else
			state.AI.ints["sightRadius"] = humanstats["sightRadius"]
		end
		
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage",
			state.renderHandle, "firstname", state.AI.strs["firstName"])
		
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage",
			state.renderHandle, "lastname", state.AI.strs["lastName"])
		

		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage",
			state.renderHandle, "health", "Excellent Health")
		
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage",
			state.renderHandle, "stamina", "Excellent Stamina")
		
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage",
			state.renderHandle, "occupation", state.AI.strs["citizenClass"])
		
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage",
			state.renderHandle, "socialClass", entityData.socialClass)
		
		if entityData.socialClass == "middle" then
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage",
					state.renderHandle, "carpentry", state.AI.skills.carpentry)
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", 
					SELF.id, "carpentry" .. "Events", state.AI.skillEvents["carpentry"])
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", 
					SELF.id, "carpentry" .. "EventsRequired", humanstats.numSkillEventsForLevel["carpentry"][state.AI.skills["carpentry"]])

			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage",
					state.renderHandle, "naturalism", state.AI.skills.naturalism)
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", 
					SELF.id, "naturalism" .. "Events", state.AI.skillEvents["naturalism"])
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", 
					SELF.id, "naturalism" .. "EventsRequired", humanstats.numSkillEventsForLevel["naturalism"][state.AI.skills["naturalism"]])

			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage",
					state.renderHandle, "smithing", state.AI.skills.smithing)
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", 
					SELF.id, "smithing" .. "Events", state.AI.skillEvents["smithing"])
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", 
					SELF.id, "smithing" .. "EventsRequired", humanstats.numSkillEventsForLevel["smithing"][state.AI.skills["smithing"]])

			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage",
					state.renderHandle, "stoneworking", state.AI.skills.stoneworking)
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", 
					SELF.id, "stoneworking" .. "Events", state.AI.skillEvents["stoneworking"])
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", 
					SELF.id, "stoneworking" .. "EventsRequired", humanstats.numSkillEventsForLevel["stoneworking"][state.AI.skills["stoneworking"]])

			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage",
					state.renderHandle, "cooking", state.AI.skills.cooking)
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", 
					SELF.id, "cooking" .. "Events", state.AI.skillEvents["cooking"])
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", 
					SELF.id, "cooking" .. "EventsRequired", humanstats.numSkillEventsForLevel["cooking"][state.AI.skills["cooking"]])


			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage",
					state.renderHandle, "science", state.AI.skills.science)
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", 
					SELF.id, "science" .. "Events", state.AI.skillEvents["science"])
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", 
					SELF.id, "science" .. "EventsRequired", humanstats.numSkillEventsForLevel["science"][state.AI.skills["science"]])


			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage",
					state.renderHandle, "farming", state.AI.skills.farming)
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", 
					SELF.id, "farming" .. "Events", state.AI.skillEvents["farming"])
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", 
					SELF.id, "farming" .. "EventsRequired", humanstats.numSkillEventsForLevel["farming"][state.AI.skills["farming"]])


			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage",
					state.renderHandle, "constructionRepair", state.AI.skills.constructionRepair)

			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage",
					state.renderHandle, "militarySkill", state.AI.skills.militarySkill)
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", 
					SELF.id, "militarySkill" .. "Events", state.AI.skillEvents["militarySkill"])
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", 
					SELF.id, "militarySkill" .. "EventsRequired", humanstats.numSkillEventsForLevel["militarySkill"][state.AI.skills["militarySkill"]])
		else
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage",
					state.renderHandle, "carpentry", 1)
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage",
					state.renderHandle, "naturalism", 1)
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage",
					state.renderHandle, "smithing", 1)
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage",
					state.renderHandle, "stoneworking", 1)
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage",
					state.renderHandle, "cooking", 1)
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage",
					state.renderHandle, "science", 1)
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage",
					state.renderHandle, "farming", 1)
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage",
					state.renderHandle, "constructionRepair", 1)
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage",
					state.renderHandle, "militarySkill", 1)
		end
		
		-- Default QoL values

		state.AI.ints.QoLSleepHappiness = 10
		state.AI.ints.QoLSleepDespair = 0
		state.AI.ints.QoLSleepAnger = 0
		state.AI.ints.QoLSleepFear = 0
		state.AI.strs.QoLSleepName = "Slept Decently"
		state.AI.strs.QoLSleepDescription = state.AI.strs.CAPITALIZED_POSSESSIVE .. " recent quality of sleep has been acceptable."
		state.AI.strs.QoLSleepIcon = "bed"
		state.AI.strs.QoLSleepIconSkin = "ui/thoughtIcons.xml"
		
		send("rendOdinCharacterClassHandler","odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLSleepHelp", "Will need proper shelter soon.")
		
		state.AI.ints.QoLSafetyHappiness = 10
		state.AI.ints.QoLSafetyDespair = 0
		state.AI.ints.QoLSafetyAnger = 0
		state.AI.ints.QoLSafetyFear = 0
		state.AI.strs.QoLSafetyName = "Reasonably Safe"
		state.AI.strs.QoLSafetyDescription = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " feels reasonably safe from harm. (DEV)"
		state.AI.strs.QoLSafetyIcon = "morale4"
		state.AI.strs.QoLSafetyIconSkin = "ui/thoughtIcons.xml"

		state.AI.ints.QoLHungerHappiness = 10
		state.AI.ints.QoLHungerDespair = 0
		state.AI.ints.QoLHungerAnger = 0
		state.AI.ints.QoLHungerFear = 0
		state.AI.strs.QoLHungerName = "Well Fed"
		state.AI.strs.QoLHungerDescription = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " is well fed and ready to work. (DEV)"
		state.AI.strs.QoLHungerIcon = "food_plate"
		state.AI.strs.QoLHungerIconSkin = "ui/thoughtIcons.xml"

		state.AI.ints.QoLWorkConditionsHappiness = 10
		state.AI.ints.QoLWorkConditionsDespair = 0
		state.AI.ints.QoLWorkConditionsAnger = 0
		state.AI.ints.QoLWorkConditionsFear = 0
		state.AI.strs.QoLWorkConditionsName = "Decent Work Conditions"
		state.AI.strs.QoLWorkConditionsDescription = state.AI.strs.CAPITALIZED_POSSESSIVE .. " workplace quality is acceptable."
		state.AI.strs.QoLWorkConditionsIcon = "workshops_category"
		state.AI.strs.QoLWorkConditionsIconSkin = "ui/orderIcons.xml"
		
		send("rendOdinCharacterClassHandler","odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLWorkConditionsHelp", "Eager to establish the Colony.")

		state.AI.ints.QoLCrowdingHappiness = 10
		state.AI.ints.QoLCrowdingDespair = 0
		state.AI.ints.QoLCrowdingAnger = 0
		state.AI.ints.QoLCrowdingFear = 0
		state.AI.strs.QoLCrowdingName = "Uncrowded"
		state.AI.strs.QoLCrowdingDescription = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " is new to the Frontier and finds the colony to be uncrowded."
		state.AI.strs.QoLCrowdingIcon = "population_icon"
		state.AI.strs.QoLCrowdingIconSkin = "ui/orderIcons.xml"


		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLSleepHappiness", state.AI.ints.QoLSleepHappiness)
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLSleepDespair", state.AI.ints.QoLSleepDespair)
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLSleepAnger", state.AI.ints.QoLSleepAnger)
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLSleepFear", state.AI.ints.QoLSleepFear)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLSleepName", state.AI.strs.QoLSleepName)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLSleepDescription", state.AI.strs.QoLSleepDescription)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLSleepIcon", state.AI.strs.QoLSleepIcon)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLSleepIconSkin", state.AI.strs.QoLSleepIconSkin)

		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLSafetyHappiness", state.AI.ints.QoLSafetyHappiness)
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLSafetyDespair", state.AI.ints.QoLSafetyDespair)
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLSafetyAnger", state.AI.ints.QoLSafetyAnger)
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLSafetyFear", state.AI.ints.QoLSafetyFear)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLSafetyName", state.AI.strs.QoLSafetyName)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLSafetyDescription", state.AI.strs.QoLSafetyDescription)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLSafetyIcon", state.AI.strs.QoLSafetyIcon)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLSafetyIconSkin", state.AI.strs.QoLSafetyIconSkin)

		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLHungerHappiness", state.AI.ints.QoLHungerHappiness)
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLHungerDespair", state.AI.ints.QoLHungerDespair)
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLHungerAnger", state.AI.ints.QoLHungerAnger)
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLHungerFear", state.AI.ints.QoLHungerFear)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLHungerName", state.AI.strs.QoLHungerName)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLHungerDescription", state.AI.strs.QoLHungerDescription)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLHungerIcon", state.AI.strs.QoLHungerIcon)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLHungerIconSkin", state.AI.strs.QoLHungerIconSkin)

		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLWorkConditionsHappiness", state.AI.ints.QoLWorkConditionsHappiness)
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLWorkConditionsDespair", state.AI.ints.QoLWorkConditionsDespair)
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLWorkConditionsAnger", state.AI.ints.QoLWorkConditionsAnger)
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLWorkConditionsFear", state.AI.ints.QoLWorkConditionsFear)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLWorkConditionsName", state.AI.strs.QoLWorkConditionsName)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLWorkConditionsDescription", state.AI.strs.QoLWorkConditionsDescription)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLWorkConditionsIcon", state.AI.strs.QoLWorkConditionsIcon)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLWorkConditionsIconSkin", state.AI.strs.QoLWorkConditionsIconSkin)

		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLCrowdingHappiness", state.AI.ints.QoLCrowdingHappiness)
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLCrowdingDespair", state.AI.ints.QoLCrowdingDespair)
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLCrowdingAnger", state.AI.ints.QoLCrowdingAnger)
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLCrowdingFear", state.AI.ints.QoLCrowdingFear)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLCrowdingName", state.AI.strs.QoLCrowdingName)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLCrowdingDescription", state.AI.strs.QoLCrowdingDescription)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLCrowdingIcon", state.AI.strs.QoLCrowdingIcon)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLCrowdingIconSkin", state.AI.strs.QoLCrowdingIconSkin)

		state.AI.ints.QoLCrowdingRating = 3
		state.AI.ints.QoLWorkConditionsRating = 3
		state.AI.ints.QoLHungerRating = 3
		state.AI.ints.QoLSafetyRating = 3
		state.AI.ints.QoLSleepRating = 3
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLCrowdingRating", state.AI.ints.QoLCrowdingRating)
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLWorkConditionsRating", state.AI.ints.QoLWorkConditionsRating)
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLHungerRating", state.AI.ints.QoLHungerRating)
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLSafetyRating", state.AI.ints.QoLSafetyRating)
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLSleepRating", state.AI.ints.QoLSleepRating)
		
		SELF.tags = EntityDB.HumanStats.tags

		if init.tag1 then
			local i = 1
			local doing = true
			while doing do
				if init["tag" .. i] then
					SELF.tags[ init["tag" .. i] ] = true
					i = i+1
				else
					doing = false
				end
			end
		end
		
		SELF.tags["citizen"] = true
		SELF.tags["character"] = true
		SELF.tags["conversable"] = true

		SELF.tags["phrenology_unknown"] = nil
          if isMurderer then
               SELF.tags["cultist_murderer"] = true
               SELF.tags["frontier_justice"] = true
          end
		
		if isFormerBandit then
			SELF.tags["former_bandit"] = true
			SELF.tags["temp_hostiles_dont_target"] = true
		end

		if init.socialClass then
			state.AI.strs["socialClass"] = init.socialClass
		else
			state.AI.strs["socialClass"] = entityData.socialClass
		end
		
		SELF.tags[ state.AI.strs["citizenClass"] ] = true;
		SELF.tags[ state.AI.strs["socialClass"] ] = true;
		SELF.tags[ state.AI.strs["socialClass"] .. "_class" ] = true
		
		-- for favour naturalists
		if state.AI.strs.citizenClass == "Naturalist" then
			SELF.tags.naturalism_jobs = true
		end
		
		if isFormerBandit then
			SELF.tags["bandit"] = nil -- important
			SELF.tags["Bandit"] = nil
		end

		if entityData.military == 1 then
			if entityData.socialClass == "middle" then
				if state.AI.skills.militarySkill <= 1 then
					state.AI.skills.militarySkill = 2
					send("rendOdinCharacterClassHandler",
						"odinRendererCharacterSetIntAttributeMessage",
						state.renderHandle,
						"militarySkill",
						state.AI.skills.militarySkill)
				end
			end

			SELF.tags[ "military" ] = true
			SELF.tags["military_and_vehicles"] = true		-- so steam knights etc. can get patrol
			state.AI.strs["loadout_tool"] = "firearm"
			state.AI.ints["grenades"] = state.AI.ints["grenadesMax"] -- and grenades, why not.
		else
			state.AI.ints["grenades"] = 0
		end
		if entityData.holy and entityData.holy == 1 then
			SELF.tags["holy"] = true
		end
		
		-- begin trait setup
		local traitNumbers = {}
		local function getRandomTrait()
			randomTrait = rand(1,#traitNames)
			if traitNumbers[randomTrait] or (state.AI.strs.socialClass == "lower" and allowedTraitLowerClass[randomTrait] == "no")
					or (state.AI.strs.socialClass == "upper" and allowedTraitUpperClass[randomTrait] == "no")
					or (state.AI.strs.socialClass == "middle" and allowedTraitMiddleClass[randomTrait] == "no")  then

				randomTrait = getRandomTrait()
			end
			return randomTrait
		end
		
		if isCriminal or isFormerBandit then
			state.AI.traits["Of Criminal Element"] = true
			if rand(1,3) == 1 then
				state.AI.traits["Brutish"] = true
			else
				if not SELF.tags["lower_class"] then
					local randomTrait = getRandomTrait()
					traitNumbers[randomTrait] = true
					state.AI.traits[ traitNames[ randomTrait ] ] = true
				end
			end
		elseif entityName == "Prisoner Overseer" then
			state.AI.traits["Prison Overseer"] = true
		else
			local numTraits = rand(1,3)
			
			if SELF.tags["lower_class"] then
				numTraits = 1
			end
			for i=1,numTraits do
				local randomTrait = getRandomTrait() -- = rand(1,#traitNames)
				traitNumbers[randomTrait] = true
				state.AI.traits[ traitNames[ randomTrait ] ] = true
			end

			traitNumbers = nil
		end
		
		if state.AI.traits["Doomed"] then SELF.tags.doomed = true end
		if state.AI.traits["Light Sleeper"] then SELF.tags.requiresbed = true end
		-- military gets training
		-- possible (but unlikely) for randoms to get training
		-- ... let's leave this in!

		if SELF.tags["military"] and entityData.socialClass == "middle" and state.AI.skills.militarySkill < 2 then
			state.AI.skills.militarySkill = 2
		end
		
		if state.AI.traits["Hale and Hearty"] then --Health bonus!
			state.AI.ints["healthMax"] = state.AI.ints["healthMax"] + 2
		end
		
		-- end trait setup

		send("rendUIManager", "uiAddColonist", SELF.id)

		if state.AI.strs["socialClass"] == "lower" then
			SELF.tags["workshop_jobs"] = true
			SELF.tags["work_party_jobs"] = true
			send( "rendCommandManager", "uiAddWorker", SELF)
			if state.AI.strs["citizenClass"] == "Labourer" then
				-- send("gameBlackboard", "gameObjectRequestNewOverseerMessage", SELF, 0)
			end
			if state.AI.strs["citizenClass"] == "Footsoldier" then
				send("gameBlackboard", "gameObjectRequestNewOverseerMessage", SELF, 2)
				makeSupplyText()
			end
			
			setProfessionTags()
			setOverseerWorkcrewFilters()
		end
		
		if state.AI.strs["genderstr"] == "male" then
			send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", SELF.id, 
			"capitalized_subjective_personal_pronoun", "He")
		elseif state.AI.strs["genderstr"] == "female" then
			send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", SELF.id, 
			"capitalized_subjective_personal_pronoun", "She")
		else 
			send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", SELF.id, 
			"capitalized_subjective_personal_pronoun", "They")
		end

		-- NV: send a message to ensure a new worker with a ? over their head produces the right alert.

		if state.AI.strs["socialClass"] == "lower" then
				send(SELF, "gameObjectSetOverseerMessage",
					SELF,
					nil,
					nil )
		end

		if state.AI.strs["socialClass"] == "upper" then
			setProfessionTags()
			setOverseerWorkcrewFilters()
		end

		if state.AI.strs["socialClass"] == "middle" then
			if state.AI.strs["citizenClass"] == "Overseer" or
				state.AI.strs["citizenClass"] == "Naturalist" or
				state.AI.strs["citizenClass"] == "Scientist" or
				state.AI.strs["citizenClass"] == "NCO" or
				state.AI.strs["citizenClass"] == "Artisan" or
				state.AI.strs["citizenClass"] == "Industrial Artisan" or
				state.AI.strs["citizenClass"] == "Vicar" or
                    state.AI.strs["citizenClass"] == "Trainee" or
				state.AI.strs["citizenClass"] == "Barber" or
				state.AI.strs["citizenClass"] == "Prisoner Overseer" then
                    
				
				state.AI.ints["workPartyIdleCounter"] = 0

				send("rendCommandManager", "uiAddWorker", SELF )
				
				SELF.tags["work_party_jobs"] = true

				local workCrewName = ""
				makeWorkCrewNames()
				
				if state.AI.strs["citizenClass"] == "Overseer" then
					workPartyType = 0
					workCrewName = state.AI.strs.workCrewNameCivilian
				elseif state.AI.strs["citizenClass"] == "NCO" then
					workPartyType = 2
					workCrewName = state.AI.strs.workCrewNameMilitary
				else
					workPartyType = 1
					workCrewName = state.AI.strs.workCrewNameCivilian
				end

				send("gameBlackboard", "gameObjectNewWorkPartyMessage", SELF, workCrewName, workPartyType)
				
				-- work crew filter defaults are handled in changeCharacterClass().
			else
				printl("ai_agent", "The game is now going to crash because you didn't set the overseer's profession correctly.")
			end
			send(SELF,"setWorkShift", 0, 0, 2, 1, 0, 0, 5, 1)

			-- set all tags & filters as appropriate.
			setProfessionTags()
			setOverseerWorkcrewFilters()
		end
		
		state.AI.walkTicks = 3
		state.AI.ints["subGridWalkTicks"] = state.AI.walkTicks
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterWalkTicks",
			state.renderHandle,
			state.AI.walkTicks)
		
		state.AI.strs["occupancyMap"] = 
		".-.\\".. 
		"-C-\\"..
		".-.\\"
		
		state.AI.strs["occupancyMapRotate45"] =  
		".-.\\".. 
		"-C-\\"..
		".-.\\"
		
		send( "gameSpatialDictionary",
			"registerSpatialMapString",
			SELF,
			state.AI.strs["occupancyMap"],
			state.AI.strs["occupancyMapRotate45"],
			true )
		
		-- set up vocalID for vocalizations; add id number vars (once we have them)
		if state.AI.strs["gender"] == "female" then
			state.AI.strs["vocalID"]= "High00"
		else
			if( rand(0,100) < 50 ) then 
				state.AI.strs["vocalID"] = "Low00"
			else
				state.AI.strs["vocalID"] = "Mid00"
			end
		end

		send("gameSpatialDictionary", "gameObjectClearBitfield", SELF)
		send("gameSpatialDictionary", "gameObjectClearHostileBit", SELF)
		
		-- I am a subject of the empire and player 1, so bits 0 and 4
		send("gameSpatialDictionary", "gameObjectAddBit", SELF, 0)
		send("gameSpatialDictionary", "gameObjectAddBit", SELF, 4)

		-- I am hostile to all other players.
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 1)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 2)
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 3)

		-- I am hostile to things with these bits set:
		if query("gameSession", "getSessionInt", "RepubliqueRelations")[1] <
			query("gameSession", "getSessionInt", "RepubliqueNeutralHostile")[1] then
			
			send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 5) -- Republique
		end
		
		if query("gameSession", "getSessionInt", "StahlmarkRelations")[1] <
			query("gameSession", "getSessionInt", "StahlmarkNeutralHostile")[1] then
			
			send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 6) -- Stahlmark
		end
		
		if query("gameSession", "getSessionInt", "NovorusRelations")[1] <
			query("gameSession", "getSessionInt", "NovorusNeutralHostile")[1] then
			
			send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 7) -- Novorus
		end
		
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 8) -- Carnivores
		-- 9 = herbivores
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 10) -- Fishpeople
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 11) -- Obeliskians
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 12) -- Selenians
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 13) -- Geometers
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 14) -- Bandits
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 15) -- Frontier Justice Targets
		send("gameSpatialDictionary", "gameObjectAddHostileBit", SELF, 16) -- Cultist Murder Targets

		state.personalMemories = {}
		state.AI.strs["originStory"] = ""
		state.AI.strs["currentMemoriesPhrase"] = ""
		state.AI.strs["religionPhrase"] = ""
		setReligionPhrase()
		state.AI.strs["politicalPrefPhrase"] = ""
		setPoliticalPrefPhrase()
		state.AI.strs["agePhrase"] = ""
		state.AI.ints["age"] = 0
		setAgePhrase()
		state.AI.strs["likesPhrase"] = ""
		state.AI.strs["physicalTraitsPhrase"] = ""
		setPhysicalTraitsPhrase()
		makeMoodText()
		makeFriendsText()
		createOriginStory()
		setDescriptiveParagraph()
		
		-- set filters for favour-based naturalist.
		if state.AI.strs.citizenClass == "Naturalist" and
			state.lockCharacterClass == true then
			
			local overseer = query("gameBlackboard",
							   "gameObjectGetOverseerMessage",
							   state.AI.currentWorkParty)[1]
			
			send("gameBlackboard", "gameSetWorkPartyMilitary", overseer, false)
			send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "hauling", false)
			send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "construction", false)
			send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "farming", false)
			send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "foraging", false)
			send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "mining", false)
			send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "chopping", false)
			send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "workshop", true)
			send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "hunting", false)
			send("gameBlackboard", "gameSetWorkPartyFilter", overseer, "military", false)
		end
		
		ready()
		
		-- welcome to the colonies!
		-- seed with a memory so we don't come into the game Tabula Rasa

		local seedMemories = EntityDB["SeedMemories"]
		
		if isCriminal then
			seedMemories = EntityDB["CriminalSeedMemories"]
		elseif isFormerBandit then
			makeMemory("Left a Life of Banditry", nil,nil,nil,nil )
			seedMemories = EntityDB["BanditSeedMemories1"]
			makeMemory( seedMemories["memories"][ rand( 1,#seedMemories["memories"] ) ], nil,nil,nil,nil )
			-- bandits get crappy seed memory set
			seedMemories = EntityDB["BanditSeedMemories2"]
		else
			makeMemory("Joined The Colonies",nil,nil,nil,nil)
		end
		
		makeMemory (seedMemories["memories"][ rand( 1,#seedMemories["memories"] ) ], nil,nil,nil,nil)
		
		local mundane_memories = {
			"Drank a Non-Alcoholic Drink",
			"Witnessed Combat",
			"Did a Dance",
			"Enjoyed a Reading of Poetry",
			"Enjoyed Smashing Something Up",
			"Drank a Jar of Brew",
			"Witnessed a Poetry Reading",
			"Prayed And Felt Better",
			"Had A Good Cower",
			"Let Off Some Steam",
			"Sleep Was Interrupted",}
		
		if not isCriminal or not isFormerBandit then
			for i=1,6 do
				makeMemory( mundane_memories[ rand(1, #mundane_memories) ],nil,nil,nil,nil)
			end
		end
		
		if init.hunger then
			local memoryName = "Feeling Hungry"
			if state.AI.ints.hunger >= 4 then
				memoryName ="Feeling Really Hungry"
			end
			makeMemory(memoryName,nil,nil,nil,nil)
		end
		state.AI.bools.ate_today = false
		
		if init.tiredness then
			local memoryName = "A Day Without Sleep"
			if state.AI.ints.hunger == 2 then
				memoryName = "Two Days Without Sleep"
			end
			makeMemory(memoryName,nil,nil,nil,nil)
		end
		
		if init.affliction then
			send(SELF,"createAffliction", nil, "blunt")
		end
		
		if init.tag1 == "temporary" then
			send("gameSession", "incSessionInt", "tempCharacterPopulation", 1)
			SELF.tags["temporary"] = true
		else
			if SELF.tags["lower_class"] == true then
				send("gameSession", "incSessionInt", "lowerClassPopulation", 1)
			elseif SELF.tags["middle_class"] == true then
				send("gameSession", "incSessionInt", "middleClassPopulation", 1)
			elseif SELF.tags["upper_class"] == true then
				send("gameSession", "incSessionInt", "upperClassPopulation", 1)
			end
		end

		mentalStateAggregator()
		recalcShiftLength() -- needs to happen after memories are injected

		send("gameSession", "incSessionInt", "colonyPopulation", 1)
		
		local popcount = query("gameSession","getSessionInt","colonyPopulation")[1]
		if popcount > query("gameSession","getSessionInt","highestPopulation")[1] then
			send("gameSession","setSessionInt","highestPopulation",popcount)
			send("gameSession","setSessionString","endGameString1", tostring(popcount) )
		end

		send(SELF,"updateSafetyQoL")
		
		-- record food for last three days; use this for food QoL
		-- assume that upon arrival they've had food appropriqate to social class.
		-- first entry is food value for *today*.
		
		local foodvalue = 1
		if state.AI.strs["socialClass"] == "middle" then
			foodvalue = 2
		elseif state.AI.strs["socialClass"] == "upper" then
			foodvalue = 3
		end
		
		state.foodrecord = {
			foodvalue,
			foodvalue,
			foodvalue,
			foodvalue,
		}
		state.AI.bools.ate_today = true
		
		state.drinkrecord = {
			1,1,1,1, -- tea across the board
		}
		SELF.tags.drank_today = true

		send(SELF,"updateHungerQoL")
		send(SELF,"updateCrowdingQoL")
		
		--let's give you some opinions on random people.
		local targets = query("gameSpatialDictionary", "allCharactersWithTagRequest", "citizen")
		if targets then
			if targets[1] then 
				if #targets[1] > 1 then
					local sortedList = {}
					for k,v in pairs(targets[1]) do
						if v ~= SELF then
							local tags = query(v,"getTags")[1]
							if not tags.dead then
								sortedList[#sortedList+1] = v
							end
						end
					end
					for i=1,6 do
						if #sortedList <= 0 then
							break
						end
						local pickNum = rand(1,#sortedList)
						local targetPick = sortedList[pickNum]
						
						local targetAmount = rand(-10,10)
						if state.AI.traits["Reclusive"] then
							targetAmount = rand(-1,1)
						elseif state.AI.strs["socialClass"] ~= "lower" then
							local AIBlock = query(targetPick, "getAIAttributes")[1]
							if AIBlock.strs["socialClass"] == "lower" and (not state.AI.traits["Common Mingler"]) then --If you're not a Common Mingler LC people are largely beneath your notice
								targetAmount = rand(-2,2)
							end
						end
						
						send(SELF,"changeFeelingsAbout",targetPick,targetAmount)
						table.remove(sortedList, pickNum)
					end
				end
			end
		end
		
		local profession = state.AI.strs["citizenClass"]
		local data = EntityDB[profession]
		if data.melee_weapons then
			-- set random default weapon
			local random_weapon = data.melee_weapons[ rand(1,#data.melee_weapons)]
			printl("ai_agent", state.AI.name .. " setting melee wep to: " .. random_weapon)
			send(SELF,"setWeapon","melee", random_weapon )
		end
		
		-- hax to get around phrenology.
		send(SELF, "phrenologistMessage")
		send(SELF,"refreshCharacterAlert")
	>>

	receive FrontierJustice()
	<<
		SELF.tags["frontier_justice"] = true
		SELF.tags["conversable"] = nil
		
		send("gameSpatialDictionary", "gameObjectAddBit", SELF, 15) -- frontier justice target
		
		send(SELF,"refreshCharacterAlert")
		
		-- abort anything we're doing to get out of Dodge
		state.AI.bools["canBeSocial"] = false
	>>
	
	receive setWorkShift( int s1, int s2, int s3, int s4, int s5, int s6, int s7, int s8)
	<<
		state.AI.ints.hour1 = s1 -- 0: enabled, 1: disabled, 2: first happiness-enabled hour ENABLED (lowest threshold) , 3: first happiness-enabled hour DISABLED (lowest threshold), 4: second happiness-enabled hour ENABLED, 5: second happiness-enabled hour DISABLED
		state.AI.ints.hour2 = s2
		state.AI.ints.hour3 = s3
		state.AI.ints.hour4 = s4
		state.AI.ints.hour5 = s5
		state.AI.ints.hour6 = s6
		state.AI.ints.hour7 = s7
		state.AI.ints.hour8 = s8

		recalcShiftLength()
	>>
	
	receive BecomeOfficeWorker(string buildingType)
	<<
		printl("ai_agent", state.AI.name .. " got becomeofficeworker : " .. tostring(buildingType ))
		
		if buildingType == "Chapel" then
			if not state.AI.traits["Spiritually Inclined"] or
				not state.AI.traits["Enthusiastic Amateur"] then
				
				makeMemory("Totally Unequipped To Be A Vicar",nil,nil,nil,nil)
			end
		else
			state.AI.ints.sermon_timer = nil
		end
		
		if buildingType then
			local buildingData = EntityDB[ buildingType ]
			if buildingData.workshifts and state.AI.strs["socialClass"] == "middle" then
				
				send(SELF,
					"setWorkShift",
					buildingData.workshifts[1],
					buildingData.workshifts[2],
					buildingData.workshifts[3],
					buildingData.workshifts[4],
					buildingData.workshifts[5],
					buildingData.workshifts[6],
					buildingData.workshifts[7],
					buildingData.workshifts[8] )

			end
			
			-- character will walk to assigned building and react to its quality
			send("gameBlackboard",
				"gameCitizenJobToMailboxMessage",
				SELF,
				nil,
				"Inspect New Workplace",
				"")
		end
	>>
	
	receive RevertOfficeWorker()
	<<
		printl("ai_agent", state.AI.name .. " receiving RevertOfficeWorker")
		
		if state.AI.strs["socialClass"] == "middle" then
			send(SELF,"setWorkShift", 0, 0, 2, 1, 0, 0, 5, 1)
		end

		if SELF.tags.military then
			SELF.tags.military = nil
			state.AI.ints.militaryTrainingPoints = 0
			--SELF.tags.militia = nil
		end
		
		changeCharacterClass("")
		
		if state.AI.curJobInstance then
			local jobs_to_cancel = {
				"Do Diplomatic Paperwork",
				"Do Event Arc Paperwork",
				"Do Event Arc Science",
				"Do Science",
				"Do Science!",
				"Administer Medical Treatment",
				"Learn Skills",
				"Do Event Arc Religion",
				"Take Confession",
				"Get Swole",
				"Get Swole (spontaneous)",
				"Peform Interrogation (event)",
				"Run To Rally Waypoint",
			}
			
			for k,v in pairs(jobs_to_cancel) do
				if state.AI.curJobInstance.name == v then
					send(SELF,"AICancelJob","unassigned from office")
					break
				end
			end
		end
	>>
	
	receive IncrementSkill(string skillType)
	<<
          local humanstats = EntityDB["HumanStats"]
		if not state.AI.skillEvents[skillType] then 
			return
		end
		if not state.AI.skills[skillType] then
			printl("ai_agent", state.AI.name .. " : WARNING: attempting to increment invalid skill: " .. tostring(skillType))
			return
		end
		
		--NOTE: IF YOU MAKE ANY CHANGES TO THE FOLLOWING "HOW MUCH SKILL AM I GETTING" BLOCK PLEASE DUPLICATE THE CHANGES IN learn_things.fsm OR THE TRAINING ACADEMY WILL FUCK UP
		if SELF.tags["worker_preached"] and SELF.tags.military == nil then --double skill gain for worker preaching
			state.AI.skillEvents[skillType] = state.AI.skillEvents[skillType] + 6
		else
			state.AI.skillEvents[skillType] = state.AI.skillEvents[skillType] + 3
		end
		
		local skillReq = humanstats.numSkillEventsForLevel[skillType][state.AI.skills[skillType]]
		
		--This is kinda brute force, but whatever.
		if state.AI.traits["Scholarly"] then --Bonus to science gain, penalty to everything else.
			if skillType == "science" then
				state.AI.skillEvents[skillType] = state.AI.skillEvents[skillType] + 1
			else
				state.AI.skillEvents[skillType] = state.AI.skillEvents[skillType] - 1
			end
		end
		if state.AI.traits["Interest in Exotic Wildernesses"] and skillType == "naturalism" then --bonus to naturalism gain
			state.AI.skillEvents[skillType] = state.AI.skillEvents[skillType] + 1
		end
		if state.AI.traits["Woodtouch"] and skillType == "carpentry" then --Bonus to Carpentry
			state.AI.skillEvents[skillType] = state.AI.skillEvents[skillType] + 1
		end
		if state.AI.traits["Stonesense"] and skillType == "stoneworking" then --Bonus to Stoneworking
			state.AI.skillEvents[skillType] = state.AI.skillEvents[skillType] + 1
		end
		if state.AI.traits["Ironborn"] and skillType == "smithing" then --Bonus to Metalworking. 
			state.AI.skillEvents[skillType] = state.AI.skillEvents[skillType] + 1
		end
		if state.AI.traits["Epicurean"] and skillType == "cooking" then --Bonus to Cooking. 
			state.AI.skillEvents[skillType] = state.AI.skillEvents[skillType] + 1
		end
		if state.AI.traits["Rustic Disposition"] and skillType == "farming" then --Bonus to Farming. 
			state.AI.skillEvents[skillType] = state.AI.skillEvents[skillType] + 1
		end
		if state.AI.traits["Craven"] and skillType == "militaryskill" then --Penalty to military skill
			state.AI.skillEvents[skillType] = state.AI.skillEvents[skillType] - 1
		end
		if state.AI.traits["Lazy"] then --Penalty to all skills
			state.AI.skillEvents[skillType] = state.AI.skillEvents[skillType] - 1
		end
		if state.AI.traits["Jack of All Trades"] and state.AI.skills[skillType] < 3 then --Bonus to skills below lv3
			state.AI.skillEvents[skillType] = state.AI.skillEvents[skillType] + 1
		end
		if state.AI.traits["Adaptable"] then ---augh special casessssss. This person gets a bonus to their best 3 skills.
			local count = 0
			local fail = false
			for k,v in pairs(state.AI.skills) do
				if state.AI.skills[skillType] < v then
					count = count + 1
					if count > 3 then
						fail = true
						break
					end
				end
			end
			if fail ~= true then
				state.AI.skillEvents[skillType] = state.AI.skillEvents[skillType] + 1
			end
		end
		if state.AI.traits["Highly Focused"] then --This person gets a bonus to their best skill, but penalties to everything else.
			local fail = false
			for k,v in pairs(state.AI.skills) do
				if state.AI.skills[skillType] < v then
					fail = true
					break
				end
			end
			if fail == true then
				state.AI.skillEvents[skillType] = state.AI.skillEvents[skillType] - 1
			else
				state.AI.skillEvents[skillType] = state.AI.skillEvents[skillType] + 1
			end
		end
		--end skillbonus block
		
		if (state.AI.skills[skillType] < 5 and
			state.AI.skillEvents[skillType] >= skillReq) then
			
			-- Skill leveup here!
			
			-- Do ding popup!
			send(SELF,"attemptEmote","skillup",3,true)
			makeMemory("Skilled Up",nil,nil,nil,nil)
		
			state.AI.skills[skillType] = state.AI.skills[skillType] + 1
			
			-- tell owned building about it.
			if state.AI.claimedWorkBuilding then
				send(state.AI.claimedWorkBuilding, "refreshSkillDisplay" )
			end
			
			state.AI.strs[skillType .. "Skill"] = EntityDB.HumanStats.skillLevelStrings[state.AI.skills[skillType]]

			local displaySkillType = skillType

			for i=1, #EntityDB.HumanStats.skillNameList do
				if skillType == EntityDB.HumanStats.skillNameList[i] then
					displaySkillType = EntityDB.HumanStats.skillDisplayNameList[i]
				end
			end
			
			local tickerText = state.AI.name .. " has learned to be better at ".. displaySkillType .. "! "
			send("rendCommandManager",
				"odinRendererTickerMessage",
				tickerText,
				"population_icon",
				"ui\\orderIcons.xml")
			
			send("rendCommandManager",
				"odinRendererCreateParticleSystemMessage",
				"Sparkle",
				state.AI.position.x,
				state.AI.position.y)
			
			send("rendOdinCharacterClassHandler",
				"odinRendererCharacterSetIntAttributeMessage",
				state.renderHandle,
				skillType,
				state.AI.skills[skillType])
			
               state.AI.skillEvents[skillType] = 0 --reset to zero!
		end

		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", SELF.id, skillType .. "Events", state.AI.skillEvents[skillType])
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", SELF.id, skillType .. "EventsRequired", humanstats.numSkillEventsForLevel[skillType][state.AI.skills[skillType]])
	>>

	receive medicalTreatmentMessage()
	<<
		-- FOR NOW REMOVES AFFLICTIONS
		table.remove(state.myAfflictions, rand(1,#state.myAfflictions) )
		state.AI.ints["health"] = state.AI.ints["healthMax"]
	>>
	
	receive forceTrait(string trait)
	<<
		state.AI.traits[trait] = true
	>>
	
	receive healRandomAffliction()
	<<
		state.AI.ints["health"] = state.AI.ints["healthMax"]
		if #state.myAfflictions > 0 then
			local afflict = rand(1,#state.myAfflictions)
			--local afflictName = state.myAfflictions[afflict].name
			send("rendOdinCharacterClassHandler",
				"odinRendererCharacterRemoveAffliction",
				SELF.id,
				#state.myAfflictions)
			
			state.AI.ints["numAfflictions"] = state.AI.ints.numAfflictions - 1
			
			table.remove(state.myAfflictions, afflict )	
			--return "afflictionNameResponse", afflictName
		end
	>>

	receive phrenologistMessage()
	<<
		SELF.tags["phrenology_unknown"] = nil
		for trait in pairs( state.AI.traits ) do
			send("rendOdinCharacterClassHandler", "odinRendererSetCharacterTraitMessage", state.renderHandle, trait, 1, false)
		end
	>>

	respond GetWorkPartyName()
	<<
		local overseer = query("gameBlackboard",
				"gameObjectGetOverseerMessage",
				state.AI.currentWorkParty)[1]
		if (overseer) then
			local nameResult = query ( "gameBlackboard", "gameObjectGetWorkPartyName", overseer)[1]
			return "workPartyNameResponse", nameResult
		else
			return "workPartyNameResponse", "No overseer found so workpartyname couldn't be returned"
		end
	>>

	respond GetWorkParty()
	<<
		return "workPartyResponse", state.AI.workParty
	>>

	receive InCombat()
	<<
		send("rendOdinCharacterClassHandler", "addCombatPanel", SELF.id)
		
		if not SELF.tags.military and state.AI.ints["inCombatTimer"] < 1 then
			-- make upset combat memory.
			send(SELF,"makeMemory","Civilian In Combat",nil,nil,nil,nil)
		end
		
		state.AI.ints["inCombatTimer"] = 10
		state.AI.ints.lastDayInCombat = query("gameSession","getSessionInt","dayCount")[1]
	>>

	receive GameObjectPlace( int x, int y ) 
	<<
		if state.AI.bools.first_placement then
			if SELF.tags["former_bandit"] then

				send("rendCommandManager",
					"odinRendererStubMessage",
					"ui\\orderIcons.xml",
					"bandit",
					state.AI.name .. " has left behind a life of Banditry!", -- header text
					state.AI.name .. " has joined your colony!", -- text description
					"Left-click to zoom. Right-click to dismiss.", -- action string
					"banditJoinedColony", -- alert type (for stacking)
					"ui//eventart//bandits.png", -- imagename for bg
					"high", -- importance: low / high / critical
					state.renderHandle, -- object ID
					30000, -- duration in ms
					0, -- "snooze" time if triggered multiple times in rapid succession
					nullHandle) -- gameobjecthandle of director, null if none
				
			end
			state.AI.bools.first_placement = false
		end
	>>

	-- LINKAGE STUFF BEYOND THIS POINT.

	receive gameObjectSetOverseerMessage(gameObjectHandle me, gameObjectHandle overseer, gameWorkParty workParty )
	<<
		local overseer_name = "none"
		if overseer then
			overseer_name = query(overseer,"getName")[1]
		end
		
		printl("ai_agent", state.AI.name .. " got gameObjectSetOverseerMessage for overseer: " ..overseer_name )
		state.AI.currentWorkParty = workParty
		
		if state.AI.strs.socialClass == "lower" then
			changeCharacterClass("")
		end
		
		send(SELF,"refreshCharacterAlert")
		state.AI.assignment = nil
		send(SELF,"updateWorkQoL")
	>>

	receive becomeFriends(gameObjectHandle friend)
	<<
		printl("ai_agent", state.AI.name .. " got becomeFriends with " .. query(friend,"getName")[1] )
		if not friend then
			return
		elseif friend.deleted then
			return
		end

		state.AI:CleanFriendsTable()

		-- YOU MAY ONLY HAVE 3 FRIENDS. EVER.
		if #state.AI.friends < state.AI.ints["maxNumberFriends"] then
			local success = state.AI:AddFriend(friend)
            if success == true then        
				local friendResult = query(friend,"getName")
				
					if not friendResult then
					return
				end
				
				local friendName = friendResult[1]
				
				if not friendName then
					return
				end	
                    
				--local tickerText = state.AI.name .. " now considers " .. friendName .. " a friend."
				--	send("rendCommandManager",
				--			"odinRendererTickerMessage",
				--			tickerText,
				--			"happy",
				--			"ui\\thoughtIcons.xml")
				
				makeMemory("Made A Friend",nil,friendName,nil,nil)
				
				makeFriendsText()
				-- force the rest, because it isn't getting picked up by MentalStateAggregator?
				setDescriptiveParagraph()
				local parsedParagraph = parseDescription(state.descriptiveParagraph)
				send("rendOdinCharacterClassHandler",
						"odinRendererSetDescriptionParagraph",
						state.renderHandle,
						parsedParagraph)
			end	
			--break
		end		
	>>
	
	receive becomeRivals(gameObjectHandle rival)
	<<
		if not rival then
			return
		elseif rival.deleted then
			return
		end

		state.AI:CleanRivalsTable()

			
		-- YOU MAY ONLY HAVE 3 RIVALS. EVER.
		if #state.AI.rivals < state.AI.ints["maxNumberFriends"] then
			local success = state.AI:AddRival(rival);
                   
			if success == true then
				local rivalResult = query(rival,"getName")
				
					if not rivalResult then
					if VALUE_STORE["showCitizenDebugConsole"] then
						printl("ai_agent", "WARNING: got bad rivalResult when querying for friend's name")
					end
					return
				end
				
				local rivalName = rivalResult[1]
				
				if not rivalName then
					if VALUE_STORE["showCitizenDebugConsole"] then
						printl("ai_agent", "WARNING: got bad rivalName when accessing rivalResult")
					end
					return
				end

                    
				send(SELF,
					"makeMemory",
					"Made An Enemy",
					nil, rivalName, nil, nil)
			end
		end
	>>

	receive detectKilled(gameObjectHandle corpse, string damageType)
	<<
		-- Name is misleading; This is not something that can really fit into hearExclamation
		-- It's for the person killing to create memories of doing the killing
		
		local memory = nil
		local memoryDescription = nil
		local otherName = false
		local otherObject = false
		local otherObjectKey = false
		
		--if you're murdering, time to stop murdering for now
		state.AI.ints.murderpoints = 0
		
		local corpseTags = query(corpse,"getTags")
		if not corpseTags or not corpseTags[1] then
			return
		else
			corpseTags = corpseTags[1]
		end
		
		if not corpseTags.dead_animal and
			not corpseTags.animal and
			not corpseTags.vermin and
			not corpseTags.dead_vermin then
			
			otherName = query(corpse, "getName")[1]
			otherObject = corpse
			otherObjectKey = "corpse"
			
			local isJustice = false
			local isCitizen = false

			if corpseTags["frontier_justice"] then
				isJustice = true
			end
			if corpseTags.citizen then
				isCitizen = true
			end
			
			local friend = false
			if corpseTags.human then
				for i=1,#state.AI.friends do
					if (state.AI.friends[i] == corpse) then
						friend = true
					end
				end
			end
			
			if not SELF.tags["killed_another"] then
				if corpseTags["fishperson"] then
					otherName = "the Fishperson (" .. otherName .. ")"
				end
				
				SELF.tags["killed_another"] = true
				
				if corpseTags["bandit"] then
                         if SELF.tags["military"] then
                              memory = "Killed a Bandit"
                              memoryDescription = "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN killed for the first time, striking down the bandit OTHER_NAME."
                         else
                              memory = "Killed for the first time"
                              memoryDescription = "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN killed for the first time, striking down the bandit OTHER_NAME. CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN is scarred by the experience."
                         end
				elseif isCitizen == true and isJustice == false then
					if friend then
						memory = "Committed Intentional Murder of a Friend"
						memoryDescription = "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN kill for the first time, committing bloody murder of one of their own friends, OTHER_NAME!"
					else
						memory = "Committed Intentional Murder"
						memoryDescription = "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN killed for the first time, committing bloody murder of OTHER_NAME!"
					end
				elseif friend then
					memory = "Killed for the first time (friend)"
				elseif state.AI.traits["Brutish"] then
					local memory = "Killed for the first time and liked it"
				elseif SELF.tags["military"] then
					memory = "Killed for the first time (military)"
				else
					memory = "Killed for the first time"
				end
			
			elseif corpseTags["fishperson"] then
				memory = "Killed a fishperson"
			elseif isCitizen and not isJustice then
				memory = "Committed Intentional Murder"
			elseif isJustice == true then
				if friend then
					memory = "Forced to Frontier Justice a Friend"
				else
					if SELF.tags.military then
						memory = "Enacted Frontier Justice (military)"
					else
						memory = "Committed Intentional Murder"
					end
				end
			end
		else
			-- Didn't kill a person. An animal maybe?
			if not corpseTags["debris"] and
				not corpseTags.vermin then
				-- Let's just make a generic "Killed some wild game" memory for now.
				memory = "Killed Some Wild Game"
			end
		end
		
		if memory and memory ~= "" then
			makeMemory(memory,memoryDescription,otherName,otherObject,otherObjectKey)
		end
	>>
	
	receive AssignmentCancelledMessage(gameSimAssignmentHandle assignment)
	<<
		if state.AI.assignment == assignment then
			printl("ai_agent", state.AI.name .. " MY ASSIGNMENT WAS CANCELLED. OH NOES")
			
			send("rendInteractiveObjectClassHandler",
				"odinRendererRemoveInteraction",
				state.renderHandle,
				"Cancel Assignment")
			
			state.AI.assignment = nil
			
		elseif state.assignment then
			-- for when you're the TARGET of an assignment
			-- (Yes, this is bloody weird, and a hack, but it's True.)
			send("rendInteractiveObjectClassHandler",
				"odinRendererRemoveInteraction",
				state.renderHandle,
				"Cancel corpse orders")
          
			state.assignment = nil
			send(SELF,"resetInteractions")
			
		else
			printl("ai_agent", state.AI.name .. " Got a stale, or erroneous, assignment cancellation message.");
		end
	>>

	receive JobLostAssignmentMessage ( gameSimJobInstanceHandle jobInstance )
	<<
		printl("ai_agent", state.AI.name .. " Job lost its assignment, resetting state.AI.assignment to NIL")
		state.AI.assignment = nil
	>>
	
	receive Update()
	<<
		-- corpse stuff
		if state.AI.bools["dead"] then
			if not SELF.tags["buried"] then
				send(SELF,"corpseUpdate")
			end
			return
		end
		
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", SELF.id, "morale", state.AI.ints["morale"]);
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", SELF.id, "health", state.AI.ints["health"]);
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", SELF.id, "healthMax", state.AI.ints["healthMax"]);
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", SELF.id, "numAfflictions", state.AI.ints["numAfflictions"]);
			
		local result = false
		if SELF.tags["military"] or SELF.tags["militia"] then
			tmEnter("decision tree block (military)")		
			result = tier1JobDecisionTree()
			if not result then result = tier2JobDecisionTree() end
			if not result then result = tier3MilitaryJobDecisionTree() end
			if not result then result = tier4MilitaryJobDecisionTree() end
			if not result then result = tier5MilitaryJobDecisionTree() end
			tmLeave()		
		else
			tmEnter("decision tree block (civilian)")
			result = tier1JobDecisionTree()
			if not result then result = tier2JobDecisionTree() end
			if not result then result = tier3CivilianJobDecisionTree() end
			if not result then result = tier4CivilianJobDecisionTree() end
			tmLeave()
		end
		
		--if not result then result = lowUrgencySurvivalDecisionTree() end
		
		if not SELF.tags.frontier_justice and
			not SELF.tags.doing_murder_rampage and
			not SELF.tags.no_sight then
			
			local sight = state.AI.ints["sightRadius"]
			
			if SELF.tags.naturalism_jobs and SELF.tags.scouting then
				sight = (sight * 2) + ( query(SELF,"getEffectiveSkillLevel","naturalism")[1] * 3)
			end

			if query("gameSession","getSessionBool","isDay")[1] then
				send("gameSpatialDictionary",
					"gridExploreFogOfWar",
					state.AI.position.x,
					state.AI.position.y,
					sight )
				
			elseif SELF.tags.sleeping then
				-- making this small, but not almost-zero; boring to watch people sleep in the dark.
				send("gameSpatialDictionary",
					"gridExploreFogOfWar",
					state.AI.position.x,
					state.AI.position.y,
					6)
			else
				-- isNight
				send("gameSpatialDictionary",
					"gridExploreFogOfWar",
					state.AI.position.x,
					state.AI.position.y,
					math.ceil(sight * 0.8) )
				
			end
		end

		if state.AI.ints["excuseMeTimer"] then
			state.AI.ints["excuseMeTimer"] = state.AI.ints["excuseMeTimer"] - 1
			if state.AI.ints["excuseMeTimer"] <= 0 then
				local occupancyMap = 
					".-.\\".. 
					"-C-\\"..
					".-.\\"
				local occupancyMapRotate45 = 
					".-.\\".. 
					"-C-\\"..
					".-.\\"
				state.AI.ints["excuseMeTimer"] = nil
				send("gameSpatialDictionary", "registerSpatialMapString", SELF, occupancyMap, occupancyMapRotate45, true )
			end
		end
		
		if state.AI.ints.updateTimer % 10 == 0 then
			character_doSecondUpdate()
		end

		if state.AI.ints.updateTimer % 30 == 0 then
			character_doThreeSecondUpdate()
		end
		
		-- if something else is overriding me - i.e. I have been picked up by a obekliskian or something...
		-- moved this below the updates; chars should still get hungry/tired, but NOT take/advance jobs
		if state.AI.thinkLocked then
			return
		end
		
		if state.AI.curJobInstance == nil then
			state.AI.canTestForInterrupts = true		-- reset testing for interrupts 
			local results = query( "gameBlackboard", "gameAgentNeedsJobMessage", state.AI, SELF )

			if results.name == "gameAgentAssignedJobMessage" then 
				state.AI.curJobInstance = results[ 1 ]
				state.AI.curJobInstance.assignedCitizen = SELF
				send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "currentJob", state.AI.curJobInstance.displayName);	
				send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "currentJobCategory", state.AI.curJobInstance.filter);	
				if VALUE_STORE[ "VerboseFSM" ] then
					if VALUE_STORE["showFSMDebugConsole"] then printl("FSM", "Citizen Update #" .. tostring(SELF) .. ": received job " .. state.AI.curJobInstance.name) end

				end
			end 
		else
			local results = query( "gameBlackboard", "gameAgentTestForInterruptsMessage", state.AI, SELF )
			if results.name == "gameAgentAssignedJobMessage" then 
				results[1].assignedCitizen = SELF

				if state.AI.curJobInstance then
					
					--[[if SELF.tags["military"] then
						local tickerText = state.AI.strs["firstName"] .. " " .. state.AI.strs["lastName"] .. " interrupts job " .. state.AI.curJobInstance.displayName .. " to do " .. results[1].displayName
						send("rendCommandManager", "odinRendererTickerMessage", tickerText, "work", "ui\\thoughtIcons.xml")
					end--]]
					-- nixing interrupt log message spam; should only do this for Fishpeople and stuff, but handle that elsewhere -dgb
					
				-- kill our current FSM and job

					if VALUE_STORE["showFSMDebugConsole"] then printl("FSM", "FSM: " .. state.AI.strs["firstName"] .. " " .. state.AI.strs["lastName"] .. " is attempting to abort job " .. state.AI.curJobInstance.displayName .. " due to an interrupt (" .. results[1].name .. ") !" ) end

					--printl("FSM index: " .. state.AI.FSMindex )

					if state.AI.FSMindex > 0 then

						-- run the abort state
						local tag = state.AI.curJobInstance:findTag( state.AI.curJobInstance:FSMAtIndex( state.AI.FSMindex  ) )
						local name = state.AI.curJobInstance:findName ( state.AI.curJobInstance:FSMAtIndex( state.AI.FSMindex ) )

						-- load up our fsm 
						local FSMRef = state.AI.curJobInstance:FSMAtIndex(state.AI.FSMindex)

						local targetFSM
						if FSMRef:isFSMDisabled() then
							targetFSM = ErrorFSMs[ FSMRef:getErrorFSM() ]
						else 
							targetFSM = FSMs[ FSMRef:getFSM() ]
						end

						local ok
						local nextState
						ok, errorState = pcall( function() targetFSM[ "abort" ](state, tag, name) end )
	
						if not ok then 
							print("ERROR: " .. errorState )
							FSM.stateError( state )
						end
					end

					state.AI.curJobInstance:abort( "Interrupt hit." )
					state.AI.curJobInstance = nil
				end

				state.AI.abortJob = true
				if reason ~= nil and reason ~= "" then
					state.AI.abortJobMessage = reason 
				end 
				if state.AI.abortJobMessage == nil then
					state.AI.abortJobMessage = ""
				end

				-- Reset the counter for next time
				state.AI.FSMindex = 0

				state.AI.curJobInstance = results[ 1 ]
				send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "currentJob", state.AI.curJobInstance.displayName);	
				send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "currentJobCategory", state.AI.curJobInstance.filter);	
			else
				if VALUE_STORE[ "VerboseFSM" ] then
					if VALUE_STORE["showFSMDebugConsole"] then printl("FSM", "Citizen Update #" .. tostring(SELF) .. ": doing job " .. state.AI.curJobInstance.name) end

				end

				local keepStepping = true
				while keepStepping do
					keepStepping = FSM.step( state ) 
				end		
			end
		end 
	>>

	receive hearExclamation( string name, gameObjectHandle exclaimer, gameObjectHandle subject )
	<<
 		if SELF.tags.dead then return end -- you're dead.
		if exclaimer == SELF then return end -- you can't hear your own exclamations.
		if state.AI.bools.asleep and state.AI.traits["Heavy Sleeper"] and (name ~= "planetaryAlignment") then return end --If you're a heavy sleeper you can't tell what's going on!!! EVER!!
		
		local memory = nil
		local memoryDescription = nil
		local otherName = false
		local otherObject = false
		local otherObjectKey = false
		
		-- This actually fits most cases.
		-- Override per-case if needed.
		if subject and not subject.deleted then
			otherName = query(subject,"getName")
			if otherName and otherName[1] then
				otherName = otherName[1]
			end
		elseif exclaimer and not exclaimer.deleted then
			otherName = query(exclaimer,"getName")
			if otherName and otherName[1] then
				otherName = otherName[1]
			end
		end
		
		if name == "witness_punching" and
			subject ~= SELF then
			
			memory = "Saw Someone Get Punched"
			
		elseif name == "descent_into_madness" then
			
			memory = "Witnessed Descent Into Madness"
			
			send("rendCommandManager",
					"odinRendererCreateParticleSystemMessage",
					"QuagSmokePuff",
					state.AI.position.x,
					state.AI.position.y )
			
			
		elseif name == "axe_murder" and
			exclaimer ~= SELF and
			subject ~= SELF then
			
			if state.AI.traits["Morbid"] or
				SELF.tags["doing_murder_rampage"] then
				memory = "Witnessed Murder And Liked It"
			else
				memory = "Witnessed Murder"
			end
				
		elseif name == "cannibalistic_butchery" then
			-- we could detect who is being butchered if we want from firstIgnored/subject
			if not SELF.tags["cannibal"] then
				memory = "Witnessed Butchery of a Human Corpse"
			end
		elseif name == "sea" then
			if state.AI.traits["Fishy Behaviour"] then
				memory = "Enthralled by the Call of the Sea"
			else
				memory = "Annoyed by the Call of the Sea"
			end
		elseif name == "creepyDance" then
			if not SELF.tags.cultist or state.AI.ints.despair < 25 then
				memory = "Witnessed Creepy Dance (maddened)"
			else
				memory = "Witnessed Creepy Dance"
			end
			
		elseif name == "heardEndTimesPreach" then
			if not state.AI.traits["Doomed"] or state.AI.ints.despair < 20 then
				memory = "Heard Someone Preaching the End Times"
			end
			SELF.tags["listening_to_sermon"] = nil
			
		elseif name == "heardCommunistAgitation" then
			if state.AI.traits["Communist"] then
				memory = "Inspired By A Comrade"
			elseif state.AI.strs.mood ~= "happy" and SELF.tags.lower_class then
				-- (state.AI.ints.anger > 20 and SELF.tags.lower_class) then
				memory = "Listened to Agitator"
			else
				memory = "Offended By Communist Agitator"
			end
			
		elseif name == "heardCogPrayer" then
			memory = "Heard Uplifting Prayer Service"
			SELF.tags["listening_to_sermon"] = nil
		elseif name == "heardMilitaryPrayer" then
				if state.AI.ints.preachShift ~= nil then --clear out any previous preach buffs
					state.AI.ints.preachShift = nil
					SELF.tags["military_preached"] = nil
					SELF.tags["worker_preached"] = nil
					SELF.tags["UC_preached"] = nil
				end
				if SELF.tags.military then
					local currentShift = query("gameSession", "getSessionInt", "currentShift")[1]
					SELF.tags["military_preached"] = true
					state.AI.traits["Preach_Motivated"] = true
					state.AI.ints.preachShift = currentShift
					memory = "Heard Military Prayer Service"
				else
					memory = "Heard Tedious Military Service"
				end
				SELF.tags["listening_to_sermon"] = nil
		elseif name == "heardWorkerPrayer" then
				if state.AI.ints.preachShift ~= nil then --clear out any previous preach buffs
					state.AI.ints.preachShift = nil
					SELF.tags["military_preached"] = nil
					SELF.tags["worker_preached"] = nil
					SELF.tags["UC_preached"] = nil
				end
				if not SELF.tags.military then
					local currentShift = query("gameSession", "getSessionInt", "currentShift")[1]
					SELF.tags["worker_preached"] = true
					state.AI.traits["Preach_Motivated"] = true
					state.AI.ints.preachShift = currentShift
					memory = "Heard Worker Prayer Service"
				else
					memory = "Heard Bland Worker Service"
				end
				SELF.tags["listening_to_sermon"] = nil
		elseif name == "heardUpperclassPrayer" then
				if state.AI.ints.preachShift ~= nil then --clear out any previous preach buffs
					state.AI.ints.preachShift = nil
					SELF.tags["military_preached"] = nil
					SELF.tags["worker_preached"] = nil
					SELF.tags["UC_preached"] = nil
				end
				if SELF.tags.science_jobs == true or SELF.tags.diplomacy_jobs == true or SELF.tags.upper_class == true or SELF.tags.academy_jobs == true then
					local currentShift = query("gameSession", "getSessionInt", "currentShift")[1]
					SELF.tags["UC_preached"] = true
					state.AI.traits["Preach_Motivated"] = true
					state.AI.ints.preachShift = currentShift
					memory = "Heard UC Prayer Service"
				else
					memory = "Heard Boring UC Service"
				end
				SELF.tags["listening_to_sermon"] = nil
		elseif name == "heardBadPrayer" then
				memory = "Heard A Very Strange Prayer Service"
				SELF.tags["listening_to_sermon"] = nil
		elseif name == "joyFreakout" then
			memory = "Witnessed Happy Freakout"
			send(SELF,"changeFeelingsAbout",exclaimer,-1)
			
		elseif name == "talkToAnimal" then
			memory = "Witnessed Talking To Animal"
			send(SELF,"changeFeelingsAbout",exclaimer,-1)
			
		elseif name == "talkToCorpse" then
			memory = "Witnessed Talking To Corpse"
			send(SELF,"changeFeelingsAbout",exclaimer,-1)
			
		elseif name == "cannibalization" and not SELF.tags.cannibal then
			memory = "Witnessed Consumption of Human Meat"
			send(SELF,"changeFeelingsAbout",exclaimer,-1)
			
		elseif name == "building_built" then
			if state.AI.traits["Patriotic"] or
				state.AI.traits["Industrious"] or
				state.AI.traits["Pioneering Spirit"] then
				
				memory = "Witnessed Construction"
			end
--[[
		elseif name == "militaryInspire" and not SELF.tags["military"] then
			local results = query( "gameHistoryDB", 
				"getRandomHistoryFragment", 
				"",
				{},
				{inspire = "yes"},
				{})

			if results.name == "noFragmentFound" then 
				send(SELF,"attemptEmote","soldiers",8,false) 
				memory = "Inspired by the Military"
			else
				return
			end
	]]
		elseif name == "black_magic" then
			if SELF.tags.cultist or state.AI.ints.despair > 50 then
				memory = "Witnessed Black Magic and Approved"
			else
				memory = "Witnessed Black Magic"
			end
		
		elseif name == "assault_witnessed" and
			exclaimer ~= SELF and
			subject ~= SELF then
			
			local subjectTags = query(subject,"getTags")[1] 
			
			if not subjectTags.frontier_justice or
				not subjectTags.marked_for_beating then
				
				send(subject,"FrontierJustice")
			end
			memory = "Witnessed Fighting Between Colonists"
			
		elseif name == "murder_witnessed" and
			exclaimer ~= SELF and
			subject ~= SELF then
			
			-- exclaimer = victim
			-- subject = murderer
			
			local victim_tags = query(exclaimer,"getTags")[1]
			if not victim_tags.murder_reported then
				send(exclaimer,"addTag","murder_reported")
				
				local alertstring = state.AI.name .. " witnessed " ..
					query(subject,"getName")[1] .. "'s murder of " ..
					query(exclaimer,"getName")[1] .. "!"
				
				send("rendCommandManager",
					"odinRendererTickerMessage",
					alertstring,
					"act_of_murder",
					"thoughtIcons.xml")
				
				send("rendCommandManager",
					"odinRendererStubMessage",
					"ui\\thoughtIcons.xml",
					"act_of_murder",
					"Foul Murder of " .. state.AI.name, -- header text
					alertstring, -- text description
					"Left-click to zoom. Right-click to dismiss.", -- action string
					"characterDeath", -- alert type (for stacking)
					"ui//eventart//murderMostFoul.png", -- imagename for bg
					"low", -- importance: low / high / critical
					state.renderHandle, -- object ID
					60000, -- duration in ms
					0, -- "snooze" time if triggered multiple times in rapid succession
					nullHandle) -- gameobjecthandle of director, null if none
			end
			
		elseif name == "black_magic_murder" and
			exclaimer ~= SELF and
			subject ~= SELF then
			
			if SELF.tags.cultist or state.AI.ints.despair > 70 then
				memory = "Witnessed Black Magic and Approved"
			else
				memory = "Witnessed Eldritch Murder"
				
				local eventQ = query("gameSimEventManager",
							"startEvent",
							"magical_murder_witnessed",
							{},
							{} )[1] 
				
				send(eventQ,"registerSubject",subject)
				send(eventQ,"registerTarget",exclaimer)
			end
			
		elseif name == "detectCombat" then
			if not SELF.tags["sleeping"] then
				if (not query(exclaimer, "getTags")[1]["animal"]) and
					subject ~= SELF then
					
					-- create Witnessed Combat memory if you don't have one already
					local results = query("gameHistoryDB",
									  "getRandomHistoryFragment",
									  "",
									  {character = SELF, },
									  {combat = "yes"},
									  {})
					
					if results.name == "noFragmentFound" then 
						memory = "Witnessed Combat"
					else
						return
					end
				end
			end
			
		elseif name == "fishperson_butcher_human" then
			if state.AI.traits["Morbid"] then -- neat!
				memory = "Witnessed Fishperson Butcher a Human Corpse (morbid)"
			else
				memory = "Witnessed Fishperson Butcher a Human Corpse"
			end
			
			-- if not hostile w/ fishpeople, trigger crisis
			--[[local hostile = query("gameSession","getSessionBool","fishpeoplePolicyHostile")[1]
			local crisis = query("gameSession","getSessionBool","fishpeopleEventActive")[1]
			
			if not hostile and not crisis then
				send("gameSession", "setSessionBool", "fishpeopleEventActive", true)
				local eventQ = query("gameSimEventManager",
								"startEvent",
								"fishpeople_crisis_butcher_human",
								{},
								{} )[1] 
						
				send(eventQ,"registerSubject",subject)
				send(eventQ,"registerTarget",exclaimer)
			end--]]
			
		elseif name == "upperClassSympathy" and state.AI.strs["socialClass"] ~= "upper" then
			if state.AI.traits["Obsequious Bootlicker"] then
				memory ="Felt Intense Joy at Being Understood by The Aristocracy"
			elseif state.AI.traits["Communist"] then
				memory = "Felt Fumbling Condescension of The Aristocracy"
			else
				memory = "Felt Understood by The Aristocracy"
			end
			
		elseif name == "detectFishpersonButcher" and	not SELF.tags.cannibal then
			
			-- the exclaimer is the fishperson butchering a corpse
			if state.AI.traits["Fishy Behaviour"] then
				
				memory = "Witnessed Butchery of a Fishperson Corpse (fishy)"
			else
				memory ="Witnessed Butchery of a Fishperson Corpse"
			end
			
		elseif name == "detectViolentDeath" and subject ~= SELF then
			if subject then
				-- the exclaimer is the one who died.
				-- the subject is the one who killed the exclaimer
	
				local wasMurder = false
				local killerTags = query(subject,"getTags")[1]
				
				if killerTags.citizen and
					not killerTags.executing_justice then
					
					wasMurder = true
				end
				
				local friend = false
				
				otherName = query(exclaimer,"getName")[1]
				otherObject = exclaimer
				otherObjectKey = "corpse"
				
				for i=1,#state.AI.friends do
					if state.AI.friends[i] == exclaimer then
						friend = true
					end
				end
				
				if friend and wasMurder then 
					memory = "Watched a Friend Get Murdered"
				elseif friend then
					memory = "Watched a Friend Die"
				else
					if state.AI.traits.Morbid then
						memory = "Witnessed Violent Death And Liked It"
					else
						memory = "Witnessed Violent Death"
					end
				end
			else
				-- violent death w/o murderer
				if state.AI.traits.Morbid then
					memory = "Witnessed Violent Death And Liked It"
				else
					memory = "Witnessed Violent Death"
				end
				
				otherName = query(exclaimer,"getName")[1]
				otherObject = exclaimer
				otherObjectKey = "corpse"
			end
			
		elseif name == "detectCriminalActivity" and
			exclaimer ~= SELF and
			not state.AI.traits["Of Criminal Element"] then
				
			memory = "Upset At Nearby Criminal Endeavours"
			
		elseif name == "fishpersonIntimidation" then
			
			if state.AI.traits["Coward"] or
				state.AI.ints["fear"] > 45 or
				state.AI.strs.mood == "afraid" then
				
				-- ahhh! run away :o
				memory = "Fishpeople Menacing"
				send("gameBlackboard",
					"gameCitizenJobToMailboxMessage",
					SELF,
					exclaimer,
					"Flee from Fishperson (skittish)",
					"enemy")

				send(SELF,"attemptEmote","fishperson_angry",6,true)
				--send(SELF,"AICancelJob","scared by fishperson")
				
			elseif (state.AI.strs.socialClass == "middle" and state.AI.skills.militarySkill > 1) 
					or (state.AI.strs.socialClass == "lower" and not SELF.tags.militia) or
					state.AI.traits["Brutish"] or
					state.AI.traits["Foolishly Brave"] or
					(state.AI.strs["mood"] == "angry" and state.AI.ints["fear"] < 45 ) then
				
				-- I'm a tough guy and I find that really irritating. >:| 
				memory = "Fishpeople Intimidation Angry"
				-- TODO push job ?
				
			elseif state.AI.traits["Fishy Behaviour"] then
				-- I love fishpeople why do they hate me so :(
				memory = "Fishpeople Intimidation Sad"
			else
				memory = "Fishpeople Inimidation Unpleasant"
				-- TODO push job ? 
			end
			
		elseif name == "detectCorpse" and exclaimer ~= SELF then
			-- exclaimer is the corpse 
			-- TODO: we should check to see if SELF has any relationship events with exclaimer.
			otherName = query(exclaimer,"getName")[1]
			otherObject = exclaimer
			otherObjectKey = "corpse"
			
			local otherTags = query(exclaimer,"getTags")[1]
			
			-- you can only be horrified by a fishperson corpse the once.
			local results = nil
			
			if otherTags.fishperson then
				results = query("gameHistoryDB", 
								"getRandomHistoryFragment", 
								"Horrified By Fishperson Corpse",
								{ character = SELF},
								{},
								{} )
			elseif otherTags.corpse and otherTags.citizen then
				if otherTags.rotted then
					results = query( "gameHistoryDB", 
									"getRandomHistoryFragment", 
									"",
									{ character = SELF, rottedCorpse = exclaimer },
									{},
									{} )
				else 
					results = query( "gameHistoryDB",
									"getRandomHistoryFragment", 
									"",
									{ character = SELF, corpse = exclaimer },
									{},
									{} )
				end
			end

			if results and results.name ~= "noFragmentFound" then
				return
			end
			
			if otherTags.fishperson then
				memory = "Horrified By Fishperson Corpse"
				otherName = query(exclaimer,"getName")[1]
				otherObject = exclaimer
				otherObjectKey = "corpse"
				
			elseif results.name == "noFragmentFound" then
				incMusic(1,10)
				incMusic(2,15)

				if isRotted then
					if (state.AI.strs.socialClass == "middle" and
						state.AI.skills.militarySkill > 1) or
						(state.AI.strs.socialClass == "lower" and
							SELF.tags.militia) or
						state.AI.traits["Adaptable"] or
						state.AI.traits["Brutish"] then
						
						-- it's fine.
						return
					else
						memory = "Saw a rotted corpse"
						otherName = query(exclaimer,"getName")[1]
						otherObject = exclaimer
						otherObjectKey = "rottedCorpse"
					end
				else
					memory = "Saw a corpse"
					otherName = query(exclaimer,"getName")[1]
					otherObject = exclaimer
					otherObjectKey = "corpse"
				end
			end
			
		elseif name == "Witnessed a Maddened Poetry Reading" or
			name == "Witnessed a Poetry Reading" or
			name == "Witnessed a Dreadful Poetry Reading" then
				
				memory = name
				
		elseif name == "frontierJusticeSeen" and exlaimer ~= SELF then
			local corpsename = query(subject, "getName" )[1]
			local friend = "false"
			
			for i=1,#state.AI.friends do			
				if (state.AI.friends[i] == subject) then
					friend = "true"
				end
			end
			
			if friend then
				memory = "Watched a Friend Die to Frontier Justice"
			else
				if state.AI.traits["Patriotic"] or
					state.AI.traits["Brutish"] or
					state.AI.traits["Morbid"] or
					state.AI.traits["Staunch Traditionalist"] then
					
					memory = "Witnessed Frontier Justice And Approved"
				else
					memory = "Witnessed the Madness of Frontier Justice"
				end
			end
		elseif EntitiesByType["emotion"][name] ~= nil then
			memory = name
		elseif name == "planetaryAlignment" then
			if SELF.tags["sleeping"] then
				memory = "Planet Dream"
			else
				memory = "Planet Wake"
			end
		elseif name == "Occult Ritual" then
			if SELF.tags.cultist or state.AI.ints.despair > 25 then
				memory = "Occult Ritual (positive)"
			else
				memory = "Occult Ritual (negative)"
			end
		end
		
		if memory and memory ~= "" then
			makeMemory(memory,memoryDescription,otherName,otherObject,otherObjectKey)
		end
	>>
	
	receive makeMemory(string memoryName,
				    string memoryDescription,
				    string otherName,
				    gameObjectHandle otherObject,
				    string otherObjectKey)
	<<
		if memoryName then
			if not SELF.tags.dead then
				makeMemory(memoryName,memoryDescription,otherName,otherObject,otherObjectKey)
			end
		end
	>>
	
	receive corpseUpdate()
	<<
		--[[if state.AI.bools["rotted"] then
			state.AI.ints["corpse_timer"] = state.AI.ints["corpse_timer"] + 1

			if state.AI.ints["corpse_timer"] % 100 == 0 then -- timer in gameticks
				results = query("gameSpatialDictionary", "allObjectsInRadiusRequest", state.AI.position, 10, true)
				if results and results[1] then
					send(results[1], "detectCorpse", SELF)
				end
			end
		else--]]
			
		if not state.AI.bools["rotted"] then
			--broadcast that there's a rotting corpse over here.
			state.AI.ints["corpse_timer"] = state.AI.ints["corpse_timer"] - 1
			
			if state.AI.ints["corpse_timer"] % 100 == 0 then -- timer in gameticks
				local results = query("gameSpatialDictionary",
								  "allObjectsInRadiusRequest",
								  state.AI.position, 10, true)
				
				if results and results[1] then
					send(results[1],"hearExclamation","detectCorpse",SELF,nil)
				end
				
				if state.AI.ints["corpse_timer"] < state.AI.ints["corpse_vermin_spawn_time_start"] and
					state.numVerminSpawned < 8 then
					
					if rand(1,8) == 8 then
						local handle = query( "scriptManager",
								"scriptCreateGameObjectRequest",
								"vermin",
								{legacyString = "Tiny Beetle" } )[1]
						
						local positionResult = query("gameSpatialDictionary", "nearbyEmptyGridSquare", state.AI.position, 3)[1]
						send(handle,
							"GameObjectPlace",
							positionResult.x,
							positionResult.y  )
						
						state.numVerminSpawned = state.numVerminSpawned +1
					end
				end
			end
			
			if state.AI.ints["corpse_timer"] <= 0 then
				-- here's your skeleton model swap
				state.AI.bools["rotted"] = true
				state.AI.bools["onFire"] = false -- because we're done with that.
				SELF.tags["burning"] = false
				SELF.tags["meat_source"] = false
				
				send( "rendOdinCharacterClassHandler",
					"odinRendererSetCharacterGeometry", 
					state.renderHandle,
					"models\\character\\body\\bipedSkeleton.upm", 
					"models\\character\\heads\\headSkull.upm",
					"none",
					"biped",
					"idle_dead")
				
				send("rendCommandManager",
					"odinRendererCreateParticleSystemMessage",
					"DustPuffXtraLarge",
					state.AI.position.x,
					state.AI.position.y)
				
			end
		end
	>>

	receive acceptTask()
	<<
		if VALUE_STORE["showCitizenDebugConsole"] then printl("citizen", state.AI.name .. " acceptTask received.") end
		-- reaction to this should be dependent on specific character traits.
		for trait in pairs( state.AI.traits ) do
			printl("ai_agent","Trait responses: acceptTask", "Character has trait: " .. trait)
			if ( trait == "Organized" ) then
				printl("ai_agent","Trait responses: acceptTask", "Character has trait: " .. trait)
				local happinessImpact = 50
				descriptiveString = "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN was delighted to have the chance to organize things recently."
			
			elseif ( trait == "Professionalism" ) then
				local happinessImpact = 50
				descriptiveString = "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN had an opportunity to exude professionalism lately."
			
			elseif ( trait == "Lazy" ) then
				local happinessImpact = 0
				descriptiveString = "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN had to do something useful recently."

			elseif ( trait == "Industrious" ) then
				local happinessImpact = 60
				descriptiveString = "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN was the model of economic efficiency recently!"

			elseif ( trait == "Enthusiastic Amateur" ) then
				--TODO: test to see if the character has this skill
			
			elseif ( trait == "Absent-minded" ) then
				--TODO: lower the efficiency of work done by this person

			elseif ( trait == "Epicurean" ) then
				--TODO: if this is a cooking or eating related job, make me happy.
			end
		end
	>>

	respond conversationTopicResponse(string topicstring, gameObjectHandle talkingTo)
	<<
		--[[ be given a topic, return how the topic makes you feel.
		Options:
			- happy
			- sad
			- afraid
			- angry
			- default?
			
			Each has a corresponding memory & set of possible animations
		]]
		
		-- Collect all possible responses to this topic
		
		local topic = EntityDB[topicstring]
		local score = 0
		local responses = {}
		local totalScore = 0
		local chosenResponse = nil
		
		local humanstats = EntityDB["HumanStats"]
		local worldstats = EntityDB["WorldStats"]

		if topic.responses then
			for key,value in pairs(topic.responses) do
			
				score = value.base
				score = score + value.class_influences[state.AI.strs["socialClass"]]
			
				if value.occupation_influences then
					if value.occupation_influences[state.AI.strs["citizenClass"]] then
						score = score+value.occupation_influences[state.AI.strs["citizenClass"]]
					end
				end

				for trait in pairs( state.AI.traits ) do
					if value.trait_influences[ trait ] then
						score = score + value.trait_influences[ trait ]
					end
				end
			
				if value.special_influences then
					-- these are all pulled from citizen.edb / world.edb
					if value.special_influences.hungry and state.AI.ints["hunger"] > humanstats["hungerWarningDays"] then
						score = score + value.special_influences.hungry
					end
					if value.special_influences.starving and state.AI.ints["hunger"] > humanstats["starvationWarningDays"] then
						score = score + value.special_influences.starving
					end
					if value.special_influences.tired and state.AI.ints["tiredness"] > humanstats["tirednessWarningDays"] then
						score = score + value.special_influences.tired
					end
					if value.special_influences.exhausted and state.AI.ints["tiredness"] > humanstats["tirednessMaxDays"] then
						score = score + value.special_influences.exhausted
					end
				end
			
				if score > 0 then
					totalScore = totalScore + score
					responses[value.response] = score
				end
			end
		
			local chosenScore = rand(0, totalScore)
			local currentScore = 0
		
			for key,value in pairs(responses) do
				currentScore = currentScore + value
				if currentScore >= chosenScore then
					chosenResponse = key
					break
				end
			end
		end
		
		-- Now do consequences.
		-- First: relationship
		--local other_id = query(talkingTo,"ROHQueryRequest")[1]
		if chosenResponse then
			if chosenResponse == "happy" then
				send(SELF, "changeFeelingsAbout", talkingTo, 1)      
			elseif chosenResponse == "angry" then
				send(SELF, "changeFeelingsAbout", talkingTo, -1)      
			end
		end
		
		if not chosenResponse then chosenResponse = "default" end
		
		send(SELF, "resetEmoteTimer")
		return "conversationResponseMessage", chosenResponse
	>>
	
	receive ExcuseMeMessage()
	<<
		-- printl("excuse me!");
		-- local occupancyMap = 
		-- "...\\".. 
		-- ".C.\\"..
		-- "...\\"
		-- local occupancyMapRotate45 = 
		-- 	".-.\\".. 
		-- 	"-C-\\"..s
		-- 	".-.\\"
		-- state.AI.ints["excuseMeTimer"] = 10
		-- send( "gameSpatialDictionary", "registerSpatialMapString", SELF, occupancyMap, occupancyMapRotate45, true )
	>>

	receive MoveAllowed(gameGridPosition pos)
	<<
		state.AI.bools["moveAllowed"] = true
		state.AI.position = pos
	>>

	receive MoveDenied()
	<<
		state.AI.bools["moveAllowed"] = false
	>>

	receive gameCitizenAttendingParty ( gameParty p, string jobCategory )
	<<
		state.AI.currentParty = p
		state.AI.partyJobCategory = jobCategory

		-- NOTE: This is not always correct! What type of party is it?
		if state.AI.currentParty then
			local memory = EntitiesByType["emotion"]["Attended a Party"]
			send( "gameHistoryDB", "createHistoryFragment", 
					memory["name"], 
					memory["type"],
					{character = SELF}, 
					memory["info"],
					memory["values"]
					)
		end
	>>
	receive Vocalize(string vocalization)
	<<
		if vocalization == nil then
			return
		end
		
		if (vocalization ~= "Converse" and vocalization ~= "Anger" and vocalization ~= "Happy") then -- no conversin or anger.
			send("rendInteractiveObjectClassHandler",
				"odinRendererPlaySFXOnInteractive",
				state.renderHandle,
				state.AI.strs["vocalID"] .. vocalization)
		end
	>>

	receive HarvestMessage( gameObjectHandle harvester, gameSimJobInstanceHandle ji )
	<<
		SELF.tags["meat_source"] = nil
		
		send("rendCommandManager",
			"odinRendererCreateParticleSystemMessage",
			"BloodSplashCentered",
			state.AI.position.x,
			state.AI.position.y)
			
		local harvesterTags = query(harvester,"getTags")[1]
		
		-- Selenians consume the meat upon .. using it.
		if not harvesterTags.selenian then
			local numSteaks =  EntityDB["HumanStats"]["numButcherOutput"] -- so gross.
	
			for s=1, numSteaks do
				local results = query( "scriptManager",
							"scriptCreateGameObjectRequest",
							"item",
							{legacyString = "long_pork"} )[1]
				
				if not results then 
					return "abort"
				else 
					local range = 1
					local positionResult = query("gameSpatialDictionary",
									   "nearbyEmptyGridSquare",
									   state.AI.position,
									   range)
					
					while not positionResult[1] do
						range = range + 1
						positionResult = query("gameSpatialDictionary",
										   "nearbyEmptyGridSquare",
										   state.AI.position,
										   range)
					end
					
					if positionResult[1].onGrid then
						send( results,
							"GameObjectPlace",
							positionResult[1].x,
							positionResult[1].y  )
					else
						send(results,
							"GameObjectPlace",
							state.AI.position.x,
							state.AI.position.y  )
					end
					
					local civ = query("gameSpatialDictionary", "gridGetCivilization", state.AI.position )[1]
					if civ == 0 then
						send(results,"ClaimItem")
					else
						send(results, "ForbidItem")
					end		
				end
			end
		end
		
		-- Put some music in the dynamics.
		incMusic(3,10)
		
		--[[send("rendOdinCharacterClassHandler",
			"odinRendererDeleteCharacterMessage",
			state.renderHandle)--]]
		
		--send("gameSpatialDictionary", "gridRemoveObject", SELF)
		
		state.AI.bools["rotted"] = true
		
		send( "rendOdinCharacterClassHandler",
			"odinRendererSetCharacterGeometry", 
			state.renderHandle,
			"models\\character\\body\\bipedSkeleton.upm", 
			"models\\character\\heads\\headSkull.upm",
			"none",
			"biped",
			"idle_dead")
				
		--send("gameBlackboard","gameObjectRemoveTargetingJobs", SELF, ji) -- no more burial.
		--destroyfromjob(SELF, ji)
	>>

	receive claimWorkBuilding( gameObjectHandle building )
	<<
		if building then
			printl("ai_agent", state.AI.name .. " got claimWorkBuilding on " .. query(building,"GetBuildingFancyName")[1] )
		else
			printl("ai_agent", state.AI.name .. " got claimWorkBuilding on NIL" )
		end
		
		state.AI.claimedWorkBuilding = building
		changeCharacterClass("")
		
		local workPartyResults = query("gameBlackboard","GetWorkPartyWorkers",SELF)[1]
		send(workPartyResults, "changeClass", "")
	>>
	
	receive deathBy( gameObjectHandle damagingObject, string damageType )
	<<
		send("rendOdinCharacterClassHandler", "removeCombatPanel", SELF.id)

--[[
		send(SELF, 
				"updateQoLElement",
				"Sleep", 
				-5, 
				-5, 
				-5, 
				-5, 
				"They died", 
				"This really hurts your restfulness", 
				"happy", 
				"ui/thoughtIcons.xml")
--]]

		FSM.abort( state, "Died.")

		if SELF.tags["temporary"] == true then
			send("gameSession", "incSessionInt", "tempCitizensDead", 1)
			send("gameSession", "incSessionInt", "tempCharacterPopulation", -1)
		else
			send("gameSession", "incSessionInt", "permCitizensDead", 1)
			
			if SELF.tags["lower_class"] == true then
				send("gameSession", "incSessionInt", "lowerClassPopulation", -1)
			elseif SELF.tags["middle_class"] == true then
				send("gameSession", "incSessionInt", "middleClassPopulation", -1)
				send("gameSession", "incSessionInt", "deadOverseers", 1)
			elseif SELF.tags["upper_class"] == true then
				send("gameSession", "incSessionInt", "upperClassPopulation", -1)
			end
		end
		
		SELF.tags["meat_source"] = true
		
		local removeFlesh = false
		if damageType == "eldritch_transformation" then
			removeFlesh = true
			SELF.tags["meat_source"] = nil
		end
		
		if SELF.tags["military"] or SELF.tags["militia"] then
			send("gameSession", "incSessionInt", "militaryCount", -1)
		end
		
		send(SELF,"refreshCharacterAlert")
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle,
			"currentJob",
			"Being Dead")
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle,
			"currentJobCategory",
			"")
		
		if damageType == "starvation" then
			local inside = query("gameSpatialDictionary",
							"gridHasSpatialTag",
							state.AI.position,
							"occupiedByStructure" )[1]
			if inside then
				local results = query("gameSpatialDictionary", "allBuildingsRequest")
				local building = false
				if results and results[1] then
					if results and results[1] then
						for k,v in pairs(results[1]) do
							if v then
								if query(v,"buildingHasSquare", state.AI.position)[1] then
									printl("ai_agent", "found building! ")
									building = v
									break
								end
							end
						end
					end
				end
				
				if building then
					if not query(building,"hasDoor")[1] then
						printl("ai_agent", state.AI.name .. " died in a building with no door!")
						if not query("gameSession","getSessionBool","diedInDoorlessBuilding")[1] then
							send("gameSession","setSessionBool","diedInDoorlessBuilding",true)
							send("gameSession", "setSteamAchievement", "diedInDoorlessBuilding")
						end
					end
				end
			end
		end
		
		local damage_death_string = damageType
		if damageType == "blunt" then
			damage_death_string = "a vicious pummeling" --from " .. damagingObject.state.AI.name -- do a query for "What's your name?"
		elseif damageType == "punch" then
			damage_death_string = "a vicious pummeling" --from " .. damagingObject.state.AI.name
		elseif damageType == "shrapnel" then
			damage_death_string = "being riddled with shrapnel"
  		elseif damageType == "bullet" then
			damage_death_string = "gunfire"
		elseif damageType == "fire" then
			damage_death_string = "fire"
		elseif damageType == "slash" then
			damage_death_string = "grievous lacerations"
		elseif damageType == "piercing" then
			damage_death_string = "terrible puncture wounds"
		elseif damageType == "explosion" then
			damage_death_string = "an explosion"
		elseif damageType == "axe_murder" then
			damage_death_string = "vicious axe-murder"
		elseif damageType == "eldritch_transformation" or
			damageType == "occult_curse" or
			damageType == "mind_blast" or
			damageType == "eldritch" then
			damage_death_string = "incomprehensible eldritch energies"
		elseif damageType == "selenian_infection" then
			damage_death_string = "terrifying otherworldly spores"
		end
		

		local justiceDeath = false
		local deathWithoutJustice = false
		local wasMurder = false
		local axeMurder = false
		local cultistMurder = false

		local damagerTags = false
		if damagingObject then
			damagerTags = query(damagingObject, "getTags")
			if damagerTags then
				damagerTags = damagerTags[1]
			end
		end
		
		if damagerTags then
			if damagerTags.citizen and
				not damagerTags.executing_justice then
				
				wasMurder = true
				-- you'll get what's comin'
				state.murderer = damagingObject
				
				if not SELF.tags.frontier_justice then
					-- MURDERDEATHKILL
					-- let the witnesses know.
					local results = query("gameSpatialDictionary",
									  "allObjectsInRadiusWithTagRequest",
									  state.AI.position,
									  "citizen",
									  12,
									  true)
					
					if results then
						send(results[1],
							"hearExclamation",
							"murder_witnessed",
							SELF,
							damagingObject)
					end
				end
			end
		end
		
		if damagerTags then
			if damagerTags["doing_cultist_murder"] then
				cultistMurder = true
				-- cultists gain power from successful murder
				send("gameSession", "incSessionInt", "cultPower", 1)
			elseif damageType == "axe_murder" then
				axeMurder = true
			end
		end	
		
		if SELF.tags["frontier_justice"] and damagerTags.executing_justice then	
			-- yes, I am tagged for justice and was killed by an executor of justice.
			justiceDeath = true
			damage_death_string = "Mob Justice"
			-- but am I mob justice?
			if SELF.tags["military"] then
				damage_death_string = "Frontier Justice executed by the military"
			end
			if SELF.tags["militia"] then
				damage_death_string = "Frontier Justice executed by the colonial militia"
			end
			
			-- let the witnesses know.
			local results = query("gameSpatialDictionary",
							  "allObjectsInRadiusWithTagRequest",
							  state.AI.position,
							  "citizen",
							  10,
							  true)
			
			if results then
				send(results[1],
					"hearExclamation",
					"frontierJusticeSeen",
					SELF,
					damagingObject)
			end
			
		elseif SELF.tags.frontier_justice and
			(not damagingObject or
			 not damagerTags.executing_justice) then
			
			-- tagged for justice, but died due to other causes
			deathWithoutJustice = true
			damage_death_string = "a killing blow not caused by an official Executor of Justice." --"Escaped Justice By Dying"
			
		else
			local results = query("gameSpatialDictionary",
						 "allObjectsInRadiusWithTagRequest",
						 state.AI.position,
						 10,
						 "citizen",
						 true)
			
			if results and damagingObject ~= SELF then
				send(results[1],
					"hearExclamation",
					"detectViolentDeath",
					SELF,
					damagingObject)
			end
		end
		
		send(SELF, "Vocalize", "Dying")
		
		local alertIconName = "skull"
		local alertIconXML = "ui\\thoughtIcons.xml"
		local alertTitle = "A Tragic Death"
		if damageType == "eldritch_transformation" then
			alertTitle = "A Horrifying Transformation"
		end
			
		local s = "" -- string
		local art = "ui//eventart//graveyard.png"
		
		if justiceDeath then
			alertIconName = "frontier_justice"
			alertIconXML = "ui\\orderIcons.xml"
			alertTitle = "Justice Is Served"
			art = "ui//eventart//reload.png"
			
			if state.AI.strs["socialClass"] == "upper" then
				s = "A noble aristocrat, " .. state.AI.name .. ", has died due to " .. damage_death_string .. "!"
			elseif state.AI.strs["socialClass"] == "lower" then
				s = "A common labourer has died due to " .. damage_death_string .. "."
			else
				s = state.AI.name .. " has died due to " .. damage_death_string .. "!"
			end
			
		elseif (deathWithoutJustice and cultistMurder) then
			alertIconName = "murder_stab"
			alertIconXML = "ui\\thoughtIcons.xml"
			alertTitle = "Cultist Murder Of Justice"
			art = "ui//eventart//magic_murder.png"
			
			if state.AI.strs["socialClass"] == "upper" then
				s = "A noble aristocrat, " .. state.AI.name .. ", has died due to " .. damage_death_string .. "!"
			elseif state.AI.strs["socialClass"] == "lower" then
				s = "A common labourer has died due to " .. damage_death_string .. "."
			else
				s = state.AI.name .. " has escaped Justice by being sacrificed in an act of " .. damage_death_string .. "!"
			end
			
		elseif deathWithoutJustice then
			alertTitle = "Justice Failed"
			
			if state.AI.strs["socialClass"] == "upper" then
				s = "A noble aristocrat, " .. state.AI.name .. ", has died due to " .. damage_death_string .. "!"
			elseif state.AI.strs["socialClass"] == "lower" then
				s = "A common labourer has died due to " .. damage_death_string .. "."
			else
				s = state.AI.name .. " has escaped Justice by dying due to " .. damage_death_string .. "!"
			end

		elseif cultistMurder then
			alertIconName = "murder_stab"
			alertIconXML = "ui\\thoughtIcons.xml"
			alertTitle = "Occult Murder"
			art = "ui//eventart//magic_murder.png"
			
			if state.AI.strs["socialClass"] == "upper" then
				s = "A noble aristocrat, " .. state.AI.name .. ", was sacrificed in an act of " .. damage_death_string .. "!"
			elseif state.AI.strs["socialClass"] == "lower" then
				s = "A common labourer was sacrificed in an act of " .. damage_death_string .. "!"
			else
				s = state.AI.name .. " was sacrificed in an act of " .. damage_death_string .. "!"
			end
			
		elseif axeMurder then
			local murdererName = query( damagingObject, "getName")[1]
			alertIconName = "murder_stab"
			alertIconXML = "ui\\thoughtIcons.xml"
			alertTitle = "Violent Murder"
			art = "ui//eventart//murderMostFoul.png"
			
			if state.AI.strs["socialClass"] == "upper" then
				s = "A noble aristocrat, " .. state.AI.name .. ", was murdered by " .. murdererName .. " in an act of disturbing violence!"
			elseif state.AI.strs["socialClass"] == "lower" then
				s = "A common labourer was murdered by " .. murdererName .. " in an act of disturbing violence!"
			else
				s = state.AI.name .. " was murdered by " .. murdererName .. " in an act of disturbing violence!"
			end
			
		else
			
			if state.AI.strs["socialClass"] == "upper" then
				s = "A noble aristocrat, " .. state.AI.name .. ", has died due to " .. damage_death_string .. "!"
			elseif state.AI.strs["socialClass"] == "lower" then
				s = "A common labourer has died due to " .. damage_death_string .. "."
			else
				s = state.AI.name .. " has died due to " .. damage_death_string .. "!"
			end
			
		end
          
		send("rendCommandManager",
			"odinRendererTickerMessage",
			s,
			alertIconName,
			alertIconXML)
		
		if state.AI.strs["socialClass"] == "upper" then
			
			-- drop all relations a touch when aristocrats die
			
			send( query("gameSession","getSessiongOH", "Empire")[1],
					"changeStanding",-10, "aristocrat death")
			send( query("gameSession","getSessiongOH", "Stahlmark")[1],
					"changeStanding",-5, "aristocrat death")
			send( query("gameSession","getSessiongOH", "Republique")[1],
					"changeStanding",-2, "aristocrat death")
			send( query("gameSession","getSessiongOH", "Novorus")[1],
					"changeStanding",-5, "aristocrat death")
			
			s = s .. " Our failure to protect our Betters has harmed relations with the Empire and our standing with the Civilized Nations of the world."
		elseif state.AI.strs["socialClass"] == "middle" then
			send( query("gameSession","getSessiongOH", "Empire")[1],
					"changeStanding",-5, "overseer death")
			local standing = query("gameSession", "getSessionInt", "EmpireRelations")[1]
			s = s .. " The death of an Overseer harms our relations with the Ministry; our standing with the Empire has decreased by 5. It is now " .. standing .. "."
		else
			send( query("gameSession","getSessiongOH", "Empire")[1],
					"changeStanding",-1, "labourer death")
			local standing = query("gameSession", "getSessionInt", "EmpireRelations")[1]
			s = s .. " The Ministry disapproves of waste of resource; our standing with the Empire has decreased by 1. It is now " .. standing .. "."
		end
		
		send("rendCommandManager",
			"odinRendererStubMessage", 
			alertIconXML,
			alertIconName,
			alertTitle, -- header text
			s, -- text description
			"Left-click to zoom. Right-click to dismiss.", -- action string
			"characterDeath", -- alert type (for stacking)
			"", --art, -- "ui//eventart//magic_murder.png", -- imagename for bg
			"low", -- importance: low / high / critical
			state.renderHandle, -- object ID
			60 * 1000, -- duration in ms
			0, -- "snooze" time if triggered multiple times in rapid succession
			nil) -- gameobjecthandle of director, null if none

		printl("ai_agent", "Death Notice: " .. state.AI.name .. " has died due to " .. damage_death_string .. "!");

		-- animation time
		if damageType == "eldritch_transformation" or damageType == "axe_murder" then
			-- already did animation in job?
			if removeFlesh then
				state.AI.bools["rotted"] = true
				state.AI.bools["onFire"] = false -- because we're done with that.
				SELF.tags["burning"] = false
				SELF.tags["meat_source"] = false
				
				send( "rendOdinCharacterClassHandler",
					"odinRendererSetCharacterGeometry", 
					state.renderHandle,
					"models\\character\\body\\bipedSkeleton.upm", 
					"models\\character\\heads\\headSkull.upm",
					"none",
					"biped",
					"idle_dead")
				
				send("rendCommandManager",
					"odinRendererCreateParticleSystemMessage",
					"DustPuffXtraLarge",
					state.AI.position.x,
					state.AI.position.y)
			end
		else
			local animName = false
			local deathAnims = {
				"death",
				"death",
				"death1",
				"death2",
				"death3",
				"death_brainmelt",
				"death_choke",
				"death_shot", 
				"death_while_fleeing", -- short leap forward
				"death_on_fire", -- medium sway fall shudder
				"death_platoon", -- long kneel reach towards sky
				"death_poet", -- long dramatic kneeling one handed grab
				"death_bulletriddled", -- quick spasms
				"deathHeadfalloff" -- quick and awesome
				}

			--if rand(0,100) == 100 then
			--	animName = "deathHeadfalloff"
			--else
				animName = deathAnims[ rand(1,#deathAnims) ]
			--end
			
			if animName then
				if removeFlesh then
					state.AI.bools["rotted"] = true
					state.AI.bools["onFire"] = false -- because we're done with that.
					SELF.tags["burning"] = false
					SELF.tags["meat_source"] = false
					
					send( "rendOdinCharacterClassHandler",
						"odinRendererSetCharacterGeometry", 
						state.renderHandle,
						"models\\character\\body\\bipedSkeleton.upm", 
						"models\\character\\heads\\headSkull.upm",
						"none",
						"biped",
						animName)
					
					send("rendCommandManager",
						"odinRendererCreateParticleSystemMessage",
						"DustPuffXtraLarge",
						state.AI.position.x,
						state.AI.position.y)
				else
					send("rendOdinCharacterClassHandler",
						"odinRendererSetCharacterAnimationMessage",
						state.renderHandle,
						animName,
						false)
				end
			end
		end

		incMusic(4,100)
    
		-- replace GSD with absolutely empty map
		send( "gameSpatialDictionary", "registerSpatialMapString", SELF, "c", "c", true )		
		
		if state.AI.claimedWorkBuilding then
			send(state.AI.claimedWorkBuilding, "ownerDied", "death")
			
			local overseer = query("gameBlackboard","gameObjectGetOverseerMessage",state.AI.currentWorkParty)[1]
			
			send("gameBlackboard", "gameSetWorkPartyClaimedBuilding", overseer, nil)
			send(state.AI.claimedWorkBuilding, "setBuildingOwner", nil)
		end

		send("gameBlackboard", "gameAgentHasDiedMessage", state.AI, SELF, true)
		state.AI.currentWorkParty = nil

		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterAttributeMessage",
			state.renderHandle,
			"nameprefix",
			"The Late ")
		
		send("gameSession","incSessionInt","colonyPopulation", -1)
		
		if SELF.tags.military then
			send("gameSession","incSessionInt","militaryDeaths", 1)
			send("gameSession", "incSteamStat", "stat_military_deaths", 1)
			if query("gameSession","getSessionInt","militaryDeaths")[1] >= 50 then
				if not query("gameSession","getSessionBool","over50MilitaryDeaths")[1] then
					send("gameSession","setSessionBool","over50MilitaryDeaths",true)
				end
			end
		end
		
		send("gameSession","incSessionInt","colonyDeaths", 1)
		if query("gameSession","getSessionInt","colonyDeaths")[1] >= 50 then
			if not query("gameSession","getSessionBool","over50deaths")[1] then
				send("gameSession","setSessionBool","over50deaths",true)
				send("gameSession", "setSteamAchievement", "over50deaths")
			end
		end

		if damageType == "eldritch_transformation" then
			send("rendInteractiveObjectClassHandler",
					"odinRendererPlaySFXOnInteractive",
					state.renderHandle,
					"Flesh Crack")	
			
			send("rendOdinCharacterClassHandler",
				"odinRendererSetCharacterAnimationMessage",
				state.renderHandle,
				"idle_dead",
				false)

			send("rendOdinCharacterClassHandler",
				"odinRendererQueueCharacterAnimationMessage",
				state.renderHandle,
				"")
		end
		
		send(SELF,"resetInteractions")
		SELF.tags.attempt_autoburial = true
		
		local hostile24results = query("gameSpatialDictionary",
									"HostileInRadius",
									SELF,
									state.AI.position,
									24)[1]
			
		if not hostile24results then
			send(SELF,"HandleInteractiveMessage","Bury Corpse (player order)",nil)
		end
	>>
	
	respond getIdleAnimQueryRequest()
	<<
		local idleAnims = {"idle", "idle_alt1", "idle_alt2", "idle_alt3",}
		if state.AI.strs["socialClass"] == "lower" then
			idleAnims[#idleAnims+1] = "dust_self_off"
		end
		
		-- a basic idle is always possible.
		local animList = { idleAnims[rand(1,#idleAnims)] } 
		
		if state.AI.ints["health"] < state.AI.ints["healthMax"] * 0.8 then
			table.insert(animList,"idle_injured")
		end
		if SELF.tags["doing_murder_rampage"] then
			table.insert(animList,"idle_murderous")
		end
		if state.AI.ints["despair"] >= 50 then
			table.insert(animList,"idle_insane")
		end	
		if state.AI.ints["fear"] >= 60 then
			table.insert(animList,"idle_shellshocked")
		end
		
		if state.AI.strs.mood == "angry" then
			table.insert(animList,"idle_angry")
		elseif state.AI.strs.mood == "despair" then
			table.insert(animList,"idle_sad")
		elseif state.AI.strs.mood == "afraid" then
			table.insert(animList,"idle_sad")
		elseif state.AI.strs.mood == "happy" then
			table.insert(animList,"idle_celebrate")
		end
		
		return "idleAnimQueryResponse", animList[ rand(1,#animList) ]
	>>
	
	receive emoteThought()
	<<
		thought = ""
		-- here we decide what to emote to the world.
		
		local humanstats = EntityDB["HumanStats"]
		local worldstats = EntityDB["WorldStats"]
		
		-- combat, life, death, majorly unhealthy status effects
		if SELF.tags["fishy_state"] then
			thought = "sea"
		elseif SELF.tags["burning"] and state.AI.strs["lastThought"] ~= "i_am_on_fire" then
			--if testThoughtVsLast("i_am_on_fire") then thought = "i_am_on_fire" end
			thought = "i_am_on_fire"
		elseif SELF.tags["fleeing"] and state.AI.strs["lastThought"] ~= "retreat" then
			--if testThoughtVsLast("retreat") then thought = "retreat" end
			thought = "retreat"
		elseif SELF.tags["shot_at"] then
			thought = "shot_at"
		elseif state.AI.ints["health"] <= div(state.AI.ints["healthMax"], 2) then
			thought = "affliction" 
		elseif state.AI.ints["hunger"] >= humanstats.hungerWarningDays then -- * 10 then
			thought = "hungry"	
		elseif state.AI.ints["tiredness"] >= humanstats.tirednessWarningDays then
			thought = "bed"
		elseif SELF.tags["has_plague"] then
			thought = "sick"
		else
			thought = state.AI.strs["mood"]
		end
		
		-- OUR MOMENT OF GLORY
		
		send("rendOdinCharacterClassHandler",
			"odinRendererCharacterExpression",
			state.renderHandle,
			"thought",
			thought,
			false )
		
		state.AI.strs["lastThought"] = thought
	>>
	
	receive ConsumeFood( gameObjectHandle food )
	<<
		--[[
			foodvalue guide:
			5: mindblowing
			4: very good /exotic
			3: upper class standard
			2: mc standard
			1: lc standard
			0: sub-par
			-1: actively bad
			(worse): starving or poison
			
			epicureans like good/exotic food more than normal
		]]
		local memoryName = "Ate something"
		local memoryDescription = nil

		local otherObject = false
		local otherObjectKey = false
		
		local hasEpicurean = false
		if state.AI.traits["Epicurean"] then
			hasEpicurean = true
		end 
		
		-- going to be SUPER careful while rebuilding the memory creation logic
		-- as there was some kind of awful error buried deep inside here
		
		local otherName = query(food,"getDisplayName")[1]
		local foodTags = query(food,"getTags")[1]
		local foodRecordEntry = 0 -- default
		-- do memory setup
		-- these are defaults
		--local memoryName = "Ate something"
		--local memoryDescription = "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN ate a " .. foodName .. " recently. It was filling."
		
		if foodTags.fishperson then
			memoryName = "Ate Fishperson Flesh"
			memoryDescription = "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN ate fishperson flesh recently. Creepy."
			if state.AI.traits["Fishy Behaviour"] then
				memoryName = "Ate Fishperson Flesh (fishy)"
				memoryDescription = "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN was forced to consume fishperson flesh recently and was horrified!"
				foodRecordEntry = -1
			elseif state.AI.traits["Epicurean"] then
				memoryDescription = "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN ate an unusual OTHER_NAME recently and found it quite evocative!"
				foodRecordEntry = 2
			end
			
			if rand(1,20) == 1 then
				SELF.tags.occult_mark_fishperson = true
			end
			
		elseif foodTags.exotic_caviar then
			memoryName = "Ate Exotic Caviar"
			
			if state.AI.traits["Epicurean"] then
				memoryDescription = "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN was delighted to eat some fine Exotic Caviar!"
				foodRecordEntry = 5
			elseif state.AI.traits["Fishy Behaviour"] then
				memoryName = "Ate Exotic Caviar (fishy)"
				foodRecordEntry = -1
			else
				foodRecordEntry = 4
			end
			-- TODO: DO SOMETHING COOL HERE w/ fishyness.
			
			if rand(1,20) == 1 then
				SELF.tags.occult_mark_fishperson = true
			end
			
		elseif foodTags.human then
			if SELF.tags.fishy_state then
				memoryName = "Consumed Delicious Human Flesh"
				foodRecordEntry = 4
			else
				incMusic(5,20)
				
				-- update cannibalism stats
				send("gameSession","incSessionInt","actsOfCannibalism", 1)
				send("gameSession","setSessionString","endGameString7",
					tostring( query("gameSession","getSessionInt","actsOfCannibalism")[1]) )
				
				send("gameSession", "setSessionBool", "cannibalismHasOccurred", true)
				send("gameSession", "setSteamAchievement", "cannibalismHasOccurred")
				send(SELF, "Vocalize", "Madness")
				
				if not SELF.tags["killed_another"] then
					SELF.tags["cannibal"] = true
				end
				
				memoryName = "Horrified At Eating Human Flesh"
				
				-- cannibalism alternate memories
				if state.AI.ints["despair"] > 50 then
					memoryName = "Enjoyed Eating Human Flesh"
					foodRecordEntry = 2
				elseif state.AI.traits["Epicurean"] then
					memoryName = "Ate Human Flesh (Epicurean)"
					foodRecordEntry = 3
				else
					foodRecordEntry = -1
				end
				
				local results = query("gameSpatialDictionary",
								  "allObjectsInRadiusWithTagRequest",
								  state.AI.position,
								  8,"citizen",true)
				if results then
					if results[1] then
						send(results[1], "hearExclamation", "cannibalization", SELF, nil)
					end
				end
			end
		elseif foodTags.exotic then
			--This is for if you have exotic food that isn't specifically one of the above, somehow.
			memoryName = "Ate Exotic Food"
			foodRecordEntry = 4
		elseif foodTags.fungus and state.AI.traits["Mushroom Lover"] then
			memoryName = "Ate Delicious Mushrooms"
			foodRecordEntry = 3
		else
			-- food defaults
			if state.AI.traits["Epicurean"] then
				memoryName = "Ate something (Epicurean)"
				foodRecordEntry = 3
			end
				
			if foodTags.raw then
				
				if not foodTags.meat then
					memoryName = "Ate something nasty"
					if state.AI.traits["Epicurean"] then
						memoryName = "Ate something nasty (Epicurean)"
						foodRecordEntry = 0
					else
						foodRecordEntry = 0
					end
				elseif foodTags.meat then
					memoryName = "Ate something nasty"
					memoryDescription = "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN ate a " .. otherName .. " recently. The raw meat was not a very tasty or civilized meal."
					
					if state.AI.traits["Epicurean"] then
						memoryName = "Ate something nasty (Epicurean)"
						memoryDescription = "CAPITALIZED_SUBJECTIVE_PERSONAL_PRONOUN ate a " .. otherName .. " recently and found it appallingly undercooked, though oddly compelling in an uncivilized sort of way."
						foodRecordEntry = 0
					else
						foodRecordEntry = 0
					end
				end
				
			elseif foodTags.cooked or foodTags.trade_food then
				
				local foodQuality = 1
				if foodTags.basic_food then
					foodRecordEntry = 1
				elseif foodTags.quality_food or foodTags.trade_food then
					foodQuality = 2
					foodRecordEntry = 2
				elseif foodTags.premium_food then
					foodQuality = 3
					foodRecordEntry = 3
				else
					printl("ai_agent", "Warning: " .. state.AI.name .. " ate some food but it was not tagged with a class value!")
					foodRecordEntry = 1
				end
				
				local personRank = 0
				if state.AI.strs["socialClass"] == "lower" then
					personRank = 1
				elseif state.AI.strs["socialClass"] == "middle" then
					personRank = 2
				elseif state.AI.strs["socialClass"] == "upper" then
					personRank = 3
				end
				
				if personRank == foodQuality then
					memoryName = "Ate something tasty"
					if state.AI.traits["Epicurean"] then
						memoryName = "Ate something tasty (Epicurean)"
					end
					
				elseif personRank < foodQuality then
					memoryName = "Ate something delicious"
					if state.AI.traits["Epicurean"] then
						memoryName = "Ate something delicious (Epicurean)"
					end
					
				else -- personRank > foodQuality
					memoryName = "Ate something unsatisfying"
					if state.AI.traits["Epicurean"] then
						memoryName = "Ate something unsatisfying (Epicurean)"
					end
				end
			else
				printl("ai_agent", "Nooo! Food somehow isn't meeting any requirements!")
			end
		end
		
		if foodTags.cooked then
			state.AI.ints["hunger"] = 0
		else
			state.AI.ints["hunger"] = state.AI.ints["hunger"] - 2
			if state.AI.ints["hunger"] < 0 then
				state.AI.ints["hunger"] = 0
			end
		end

		if memoryName then
			makeMemory(memoryName,memoryDescription,otherName,otherObject,otherObjectKey)
		end
		
		if state.AI.ints["hunger"] < EntityDB["HumanStats"]["starvationWarningDays"] and
			SELF.tags["starving"] then
			
			send(SELF,"refreshCharacterAlert")
			
			-- no longer starving
			SELF.tags["starving"] = nil
			send("gameSession", "incSessionInt", "starvingCount", -1)
		end
		
		table.insert(state.foodrecord,1,foodRecordEntry)
		
		state.AI.bools.ate_today = true
		send(SELF,"updateHungerQoL")
		
	>>
	
	receive ConsumeDrink( gameObjectHandle booze )
	<<
		local tags = query(booze,"getTags")[1]
		if tags then
			if tags["cult_juice"] then
				send(SELF,"makeMemory","Drank Cult Juice",nil,nil,nil,nil )
			end
			
			if tags["anger_juice"] then
				send(SELF,"makeMemory","Drank Anger Juice",nil,nil,nil,nil )
			end
			
			if tags["happy_juice"] then
				send(SELF,"makeMemory","Drank Happy Juice",nil,nil,nil,nil )
			end
			
			if tags["spirits"] then
				makeMemory("Drank a Bottle of Spirits",nil,nil,nil,nil)
				state.AI.ints["inebriation"] = state.AI.ints["inebriation"] + 10	
			elseif tags["booze"] then
				makeMemory("Drank a Jar of Brew",nil,nil,nil,nil)
				state.AI.ints["inebriation"] = state.AI.ints["inebriation"] + 5
			elseif tags["laudanum"] then
				makeMemory("Drank a Bottle of Laudanum",nil,nil,nil,nil)
				state.AI.ints["inebriation"] = state.AI.ints["inebriation"] + 15
				
				if not state.AI.traits["Laudanum Fiend"] and
					rand(1,100) == 1 then
					
					local tickerText = "Oh no, CHARACTER_NAME has become a Laudanum Fiend due to too much self-medication!"
					send("rendCommandManager",
						"odinRendererTickerMessage",
						parseDescription(tickerText),
						"laudanum_bottle",
						"ui\\commodityIcons.xml")
					
					state.AI.traits["Laudanum Fiend"] = true
					send("rendOdinCharacterClassHandler",
						"odinRendererSetCharacterTraitMessage",
						state.renderHandle,
						"Laudanum Fiend",1,false)
				end
				
			elseif tags["sulphur_tonic"] then
				makeMemory("Drank Sulphur Tonic",nil,nil,nil,nil)
				
				state.AI.ints["inebriation"] = state.AI.ints["inebriation"] + 15
				send(SELF,"damageMessage", booze, "chemical_burn", 7, "" )
				
			elseif tags["fish_juice"] then
				makeMemory("Drank Fish Juice",nil,nil,nil,nil)
				
				send("gameBlackboard",
					"gameCitizenJobToMailboxMessage",
					SELF,    
					nil,
					"Fishy Transformation 1",
					"") 

			else
				makeMemory("Drank a Non-Alcoholic Drink",nil,nil,nil,nil)
				-- TODO: coffee, tea?
			end
		end
	>>
	
	receive spawnGibs()
	<<  
		-- TODO should pull gib stats from colonists
		--local gibTable = EntityDB["HumanStats"].gibs
		local handle = query("scriptManager",
						 "scriptCreateGameObjectRequest",
						 "clearable",
						 { legacyString = "Gibs" } )[1]

		if not handle then 
			return "abort"
		else 
			local range = 1
				positionResult = query("gameSpatialDictionary", "nearbyEmptyGridSquare", state.AI.position, range)
				while not positionResult[1] do
					range = range + 1
					positionResult = query("gameSpatialDictionary", "nearbyEmptyGridSquare", state.AI.position, range)
				end
			
			if positionResult[1].onGrid then
				send( handle, "GameObjectPlace", positionResult[1].x, positionResult[1].y  )
			else
				send( handle, "GameObjectPlace", state.AI.position.x, state.AI.position.y  )
			end
	    end
	>>
	
	receive makeNewReligionPhrase( string newPhrase )
	<<
		state.AI.strs["religionPhrase"] = newPhrase
	>>
	
	receive changeFeelingsAbout( gameObjectHandle other, int change )
	<<
		-- receive a message about someone changing feelings, update stuff & give feedback appropriately.
		-- attach this to "become friends" & etc
		-- and find out everywhere we attempt to change feelings about people.
		
		if not state.AI.feelingsAbout[ other ] then 
			state.AI.feelingsAbout[ other ] = change
		else 
			state.AI.feelingsAbout[ other ] = state.AI.feelingsAbout[ other ] + change
		end

		if state.AI.feelingsAbout[ other ] >= state.AI.ints["relationshipThresholdFriend"] then
			send(SELF, "becomeFriends", other)
			send(SELF, "Vocalize", "Happy")
			
			-- becomes harder and harder to make more friends.
			state.AI.ints["relationshipThresholdFriend"] = state.AI.ints["relationshipThresholdFriend"] * 2
		end
		
		if state.AI.feelingsAbout[ other ] <= 0 then
			
			-- end friendship if friend
			local friendkey = false
			for i=1,#state.AI.friends do
				if state.AI.friends[i] == other then
					friendkey = true
					break
				end
			end
			
			if friendkey then
				state.AI:RemoveFriend(friend);
				state.AI.ints["relationshipThresholdFriend"] = div(state.AI.ints["relationshipThresholdFriend"], 2)
				local friendName = query(other,"getName")[1] 
				
				send(SELF,"makeMemory","Lost A Friend",nil,friendName,nil,nil)
				
				makeFriendsText()
				
				-- force the rest, because it isn't getting picked up by MentalStateAggregator?
				setDescriptiveParagraph()

				local parsedParagraph = parseDescription(state.descriptiveParagraph)
				send("rendOdinCharacterClassHandler",
					"odinRendererSetDescriptionParagraph",
					state.renderHandle,
					parsedParagraph)
			end
		end
		
		if  state.AI.feelingsAbout[ other ] <= -3 then
			printl("ai_agent", state.AI.name .. " becoming rival w/ " .. query(other,"getName")[1] )
			send(SELF, "becomeRivals", other )
		end
	>>
	
	respond getFeelingsAbout( gameObjectHandle other )
	<<
		if not state.AI.feelingsAbout[ other ] then
			return "feelingsAboutResponse", 0
		else
			return "feelingsAboutResponse", state.AI.feelingsAbout[ other ]
		end
	>>

	respond getFirstName()
     <<
          return "getNameResponse", state.AI.strs["firstName"]
     >>

	 respond getLastName()
     <<
		return "getNameResponse", state.AI.strs["lastName"]          
     >>
	
	respond getProfession()
     <<
          return "getProfResponse", state.AI.strs["citizenClass"]
     >>
	
	respond getMood()
     <<
          return "getMoodResponse", state.AI.strs["mood"]
     >>
	
	respond hasTrait( string traitname )
	<<
		if state.AI.traits[ traitname ] then
			return "hasTraitResponse", true
		else
			return "hasTraitResponse", false
		end
	>>

	receive newWorkShift(int shiftNumber)
	<<
		if SELF.tags.dead then
			
			-- peform burial check here
			if not SELF.tags.buried and
				not state.assignment and
				SELF.tags.attempt_autoburial then
				
				local hostile24results = query("gameSpatialDictionary",
										"HostileInRadius",
										SELF,
										state.AI.position,
										24)[1]
				
				if not hostile24results then
					send(SELF,"HandleInteractiveMessage","Bury Corpse (player order)",nil)
				end
			end
			return
		end

		if SELF.tags.middle_class then
			recalcShiftLength()
		end

		-- Do time-based updates.
		if shiftNumber == 1 then
			send(SELF,"ToSunrise")
			
			if state.AI.ints.tiredness > 0 and state.lastSleepDay < query("gameSession","getSessionInt","dayCount")[1] then
				send(SELF,"updateSleepQuality",nil,nil)
			end
		
		elseif shiftNumber == 2 then
			-- might as well do this more often.
			send(SELF,"updateSafetyQoL")
		elseif shiftNumber == 3 then
			-- noon!
			send(SELF,"updateWorkQoL")
		elseif shiftNumber == 4 then
			if state.AI.strs["citizenClass"] == "Vicar" then
				SELF.tags["can_preach"] = true
			end
		elseif shiftNumber == 6 then
			send(SELF,"ToDusk")
			-- dusk? why not.
			send(SELF,"updateSafetyQoL")
		elseif shiftNumber == 7 then
			send(SELF,"Nightfall")
			send(SELF,"updateWorkQoL")
		elseif shiftNumber == 8 then
			send(SELF,"ToMidnight") 
		end
		
		if state.AI.ints.preachShift ~= nil then --Reset all tags that have to do with preaching, if they have expired
			local shiftTotal = 0
			if shiftNumber < state.AI.ints.preachShift then --assume a day has passed.
				shiftTotal = 8 + shiftNumber - state.AI.ints.preachShift 
			else
				shiftTotal = shiftNumber - state.AI.ints.preachShift
			end
			if shiftTotal >= 7 then
				state.AI.ints.preachShift = nil
				SELF.tags["military_preached"] = nil
				SELF.tags["worker_preached"] = nil
				SELF.tags["UC_preached"] = nil
				state.AI.traits["Preach_Motivated"] = nil
			end
		end

		if state.AI.claimedWorkBuilding then
			local isOnShift = query("gameBlackboard",
							    "gameObjectGetWorkPartyOnShift",
							    state.AI.currentWorkParty)[1]

			if isOnShift then
				send(state.AI.claimedWorkBuilding,"addTag", "overseer_active")
			else
				send(state.AI.claimedWorkBuilding,"removeTag", "overseer_active")
			end
		end

		if query("gameSession", "getSessionBool", "blockWeatherMemories")[1] ~= true then
			
			local biome = query("gameSession", "getSessionString", "biome")[1]
			if (biome == "desert" or biome == "tropical") and
				shiftNumber < 6 and
				shiftNumber > 2 and
				not state.AI.traits["Hale and Hearty"] and
				not state.AI.traits["Pioneering Spirit"] then
				 
				if rand(1,8) == 1 then
					send(SELF,"makeMemory","Overheated",nil,nil,nil,nil)
				end
				
			elseif biome == "cold" and
				not state.AI.traits["Hale and Hearty"] and
				not state.AI.traits["Pioneering Spirit"] then

				if rand(1,4) == 1 or shiftNumber > 6 then
					if not query("gameSpatialDictionary","gridHasSpatialTag",state.AI.position,"occupiedByStructure" )[1] then
						send(SELF,"makeMemory","Chilled By The Cold",nil,nil,nil,nil)
					end
				end
			end
		end
	>>
	
	receive ToDusk()
	<<
		if SELF.tags.dead then
			return
		end
		
		state.AI.ints.hunger = state.AI.ints.hunger + 1
		-- do food quality record here!
		
		if state.AI.ints["hunger"] >= state.AI.ints["starvationDeathHunger"] then
			send(SELF, "deathBy", SELF, "starvation")
			return
		end
		
		-- starvation warning 
		if state.AI.ints["hunger"] >= EntityDB["HumanStats"]["starvationWarningDays"] and
			not SELF.tags["starving"] and
			not SELF.tags.frontier_justice and
			not SELF.tags.selenian_infested then
			
			-- send player a warning, flip starvation flag
			SELF.tags["starving"] = true
			
			makeMemory("Desperately Hungry",nil,nil,nil,nil)
			
			send("rendCommandManager",
				"odinRendererStubMessage",
				"ui\\thoughtIcons.xml", -- iconskin
				"hungry", -- icon
				"Starvation", -- header text
				state.AI.name .. " is starving and will die soon without food!", -- text description
				"Left-click to zoom. Right-click to dismiss.", -- action string
				"characterStarvation", -- alert type (for stacking)
				"ui//eventart//tasty_pickles.png", -- imagename for bg
				"low", -- importance: low / high / critical
				state.renderHandle, -- object ID
				30000, -- duration in ms
				0, -- "snooze" time if triggered multiple times in rapid succession
				nil) -- gameobjecthandle of director, null if none

			
			local tickerText = state.AI.strs["firstName"] .. " " .. state.AI.strs["lastName"] .. " is starving and will die soon without food!"
			send("rendCommandManager",
				"odinRendererTickerMessage",
				tickerText,
				"hungry",
				"ui\\thoughtIcons.xml")
			
			send(SELF,"refreshCharacterAlert")
			
			send("gameSession", "incSessionInt", "starvingCount", 1)
			-- free starving but stuck colonists?
			FSM.abort( state, "starving" )
		end
		
		-- upset hunger memories
		if state.AI.ints.hunger > 1 then
			local memoryName = "Feeling Hungry"

			if state.AI.ints["hunger"] <= 4 then
				memoryName = "Feeling Really Hungry"
				
			elseif state.AI.ints["hunger"] <= 7 then
				memoryName = "Feeling Starvation"
				
			elseif state.AI.ints["hunger"] >= 8 and
				not SELF.tags.selenian_infested then
				
				local cannibalismCounter = 0
				if state.AI.ints.despair > 19 then
					cannibalismCounter = math.floor( state.AI.ints.despair * 0.1 ) + math.floor( state.AI.ints.anger * 0.1 )
				end
				local cannibalismTraits = { "Fishy Behaviour",
					"Strange",
					"Morbid",
					"Brutish",
					"Big Game Hunter"}
			
				for k,trait in pairs( state.AI.traits ) do
					for k2,v in pairs(cannibalismTraits) do
						if trait == v then
							cannibalismCounter = cannibalismCounter + 2
						end
					end
				end
				
				local isViolent = false
				if cannibalismCounter > 6 then
					isViolent = true
				end
				
				printl("ai_agent", state.AI.name .. " cannibalismCounter = " .. cannibalismCounter .. " / isViolent?: " .. tostring(isViolent))
					
				if isViolent then 
					memoryName = "Murderously Starving"
					printl("ai_agent", state.AI.name .. " is Murderously Starving!")
					if not SELF.tags.cannibalistic_murderer then
						
						send("gameBlackboard",
							"gameCitizenJobToMailboxMessage",
							SELF,	
							nil,
							"Begin Cannibalistic Murder",
							"")
						
					end
				else
					memoryName = "Hopelessly Starving"
				end
			end
			
			makeMemory(memoryName,nil,nil,nil,nil)
		end
	>>

	receive Nightfall()
	<<
		-- It's transitioning to nighttime! Do stuff you'd do at night.
		if SELF.tags["dead"] and not SELF.tags.last_rites_performed then
			
			local ghostChance = 2
			
			if state.murderer and SELF.tags.murder_avenged then
				ghostChance = ghostChance +1
			elseif state.murderer and not SELF.tags.murder_avenged then
				ghostchance = ghostChance + 8
			end
			
			if not SELF.tags.buried or
				SELF.tags.occult_mark_death then
				ghostChance = ghostChance + 5
			end
			
			if rand(1,100) < ghostChance then
					
				-- maybe spawn a ghost.
				-- ghostgoals
				-- ghostchance
				-- unburied
				-- if murdered
				-- (and killer still alive)
				
				local goal = "haunting"
				if state.murderer and
					not SELF.tags.murder_avenged then
					goal = "vengeance"
				elseif not SELF.tags.buried then
					goal = "burial"
				end
				
				local spawnTable = {legacyString = "Spectre",
								name = state.AI.name,
								goal = goal }
				
				local handle = query( "scriptManager",
									"scriptCreateGameObjectRequest",
									"spectre",
									spawnTable )[1]
							
				send(handle,
					"GameObjectPlace",
					state.AI.position.x,
					state.AI.position.y  )
				
				state.mySpectre = handle
				send(state.mySpectre,"registerOwner",SELF)
				if goal == "vengeance" and state.murderer then
					send(state.mySpectre,"registerHauntingTarget",state.murderer)
				end
			end
			return
		end
		
		-- lose invulnerability upon night.
		if SELF.tags.temp_hostiles_dont_target then
			SELF.tags.temp_hostiles_dont_target = nil
		end
		
		if not SELF.tags.selenian_infested then
			state.AI.ints.tiredness = state.AI.ints.tiredness + 1
		end
		
		
	>>
	
	receive ToSunrise()
	<<
		-- It's daytime! Do daytime stuff if you want.
		-- TODO: if sleeping outside, wake up.
		
		if state.AI.ints.hunger > 0 then --This is kinda hacky but it allows for that half-day distinction while preventing people from eating 2 food per day
			state.AI.ints.hunger = state.AI.ints.hunger + 1
		end
		
		if query("gameSession", "getSessionInt", "dayCount")[1] > 1 then
			local resultText = query("gameHistoryDB",
					"generateDailyJournalLog",
					SELF,
					state.AI)
			state.AI.strs.dailyJournalText = parseDescription(resultText[1])
			send("gameHistoryDB", "clearDailyJournal", SELF, state.AI)
		end
		if query("gameSession", "getSessionInt", "dayCount")[1] % 7 == 0 then
			local resultText = query("gameHistoryDB",
					"generateWeeklyJournalLog",
					SELF,
					state.AI)
			state.AI.strs.weeklyJournalText = parseDescription(resultText[1])
			send("gameHistoryDB", "clearWeeklyJournal", SELF, state.AI)
		end

		if not SELF.tags.dead then
			-- let's do crowding QoL here too.
			send(SELF,"updateCrowdingQoL")
			
		end
	>>
	
 
     receive damageMessage( gameObjectHandle damagingObject, string damageType, int damageAmount, string onhit_effect )
     <<
		if SELF.tags.dead then
			return
		end
		
		local damagerTags = {}
		if damagingObject then 
			damagerTags = query(damagingObject, "getTags")[1]
		end

		-- if attacked by human, justice them
		-- AHHHH CLEAN THIS UP
		if damagerTags.citizen and 
			(damagerTags.cannibalistic_murderer or
			damagerTags.cultist_murderer or
			damagerTags.doing_murder_rampage) and
			not SELF.tags.marked_for_beating then
			
			-- do witness check here.
			
			send(damagingObject,"FrontierJustice")
		end

		if damagerTags.animal then
			send(damagingObject,"makeHostile")
		end
		
		local memoryName = "Damaged"
		
		-- create upsetting memories
		if damagerTags.fishperson and
			not SELF.tags.hostile_agent and
			not damagerTags.hostile_agent and
			not damagerTags.assault_crisis_launched and
			onhit_effect ~= "no_frontier_justice" then
			
			-- hurt by fishperson!
			send(damagingObject, "makeHostile")

			if state.AI.traits["Fishy Behaviour"] then
				memoryName = "Fishpeople Menacing (fishy)"
			else
				memoryName = "Fishpeople Menacing"
			end
			
			-- launch crisis decision over fishperson assault ?
			--[[local friendly = query("gameSession", "getSessionBool", "fishpeoplePolicyFriendly")[1]
			local denial = query("gameSession","getSessionBool","fishpeoplePolicyDenial")[1]
			local hostile = query("gameSession", "getSessionBool", "fishpeoplePolicyHostile")[1]
			local crisisActive = query("gameSession", "getSessionBool", "fishpeopleCrisisActive")[1]
			
			if not crisisActive then
				if friendly == true or
					denial == true or
					(friendly == false and denial == false and hostile == false ) then

					local eventQ = query("gameSimEventManager",
									"startEvent",
									"fishpeople_crisis_assault",
									{},
									{} )[1] 
					
					send(eventQ,"registerSubject",SELF)
					send(eventQ,"registerTarget",damagingObject)
					send(damagingObject,"addTag","assault_crisis_launched")
				end
			end--]]
			
		elseif damagerTags.obeliskian and
			not SELF.tags.hostile_agent then
			
			if damageType == "eldritch" then 
				memoryName = "Suffered From Obeliskian Aetheric Attack"
			else
				memoryName = "Attacked By An Obeliskian"
			end
			
		elseif damageType == "selenian" and
			not SELF.tags.hostile_agent then
			
			memoryName = "Suffered From Selenian Field Effect"
			
		elseif damagerTags.selenian and
			not SELF.tags.hostile_agent then
			
			if damageType == "eldritch" then 
				memoryName = "Suffered Selenian Field Effect"
			else
				memoryName = "Attacked By A Selenian"
			end
			
		elseif damagerTags.citizen and
			not damagerTags.executing_justice and
			not SELF.tags.frontier_justice and
			not SELF.tags.corpse and
			onhit_effect ~= "no_frontier_justice" then
			
			--check for witnesses if this is a murder
			
			local attackerName = query(damagingObject,"getName")[1]
				send("rendCommandManager",
					"odinRendererTickerMessage",
					attackerName .. " has violently attacked " ..
						state.AI.name .. "!",
					"act_of_murder",
					"ui\\thoughtIcons.xml")
				
			send(damagingObject,"FrontierJustice")
			
			local results = query("gameSpatialDictionary",
						  "allObjectsInRadiusWithTagRequest",
						  state.AI.position,
						  12,
						  "citizen",
						  true)
			
			if results then
				send(results[1],"hearExcalamation","assault_witnessed", SELF, nil)
			end
			memoryName = "Attacked By A Murderer"
		end
		
		makeMemory(memoryName,nil,nil,nil,nil)
     >>
	
	receive despawnTemporaryCharacter()
	<<
		printl("ai_agent", "despawning temporary character: " .. state.AI.name )
		FSM.abort( state, "Despawning.")
		send(SELF,"AICancelJob", "despawning")
		send(SELF,"ForceDropEverything")
		send("gameSession", "incSessionInt", "tempCharacterPopulation", -1)
		
		-- No, don't get precious with this. Just delete self and do cleanup.
		-- doing death cleanup.
		sleep()
		
		if state.AI.claimedWorkBuilding then
			send(state.AI.claimedWorkBuilding, "ownerDied", "left")
				local overseer = query("gameBlackboard",
						"gameObjectGetOverseerMessage",
						state.AI.currentWorkParty)[1]
				send("gameBlackboard", "gameSetWorkPartyClaimedBuilding", overseer, nullHandle)
			
			send(state.AI.claimedWorkBuilding,
				"setBuildingOwner",
				nullHandle)
		end
		
		if SELF.tags["cultist"] then
			send("gameSession","incSessionInt", "cultPower", -1)
			if state.group then
				send(state.group,"removeMember", SELF, "left_map")
				if SELF.tags["cult_leader"] then
					send(state.group, "leaderDied", "left the colony")
				end
			end
		end
		
		-- hmm. Not quite appropriate.
		send("gameBlackboard", "gameAgentHasDiedMessage", state.AI, SELF, false)
		send("rendUIManager", "uiRemoveColonist", SELF.id) -- this is good.
		
		state.AI.currentWorkParty = nil
		send("gameSession","incSessionInt","colonyPopulation", -1)
		
		state.AI.bools["canBeSocial"] = false
		send("rendOdinCharacterClassHandler",
			"odinRendererHideCharacterMessage",
			state.renderHandle,
			true)

		printl("ai_agent", state.AI.name .. " being removed from game.")
		
		send("rendOdinCharacterClassHandler",
			"odinRendererDeleteCharacterMessage",
			state.renderHandle)
		
		send("gameSpatialDictionary",
			"gridRemoveObject",
			SELF)
		
		send("gameBlackboard",
			"gameObjectRemoveTargetingJobs",
			SELF,
			nil)
		
		destroyfromjob(SELF, nil)
	>>
	
	receive changeClass (string newclass)
	<<
		changeCharacterClass( newclass )
	>>

	receive recalcShiftLengthMessage ()
	<<
		recalcShiftLength()
	>>
	
	respond characterParseString( string s)
	<<
		newstring = parseDescription(s)
		return "characterParseString", newstring
	>>
	
	receive InteractiveMessage( string messagereceived )
	<<
		if messagereceived == "Frontier Justice" then
			send(SELF,"FrontierJustice")
			return
		end
		
		if messagereceived == "Spontaneous Human Combustion" then
			printl("setting myself on fire")	
			send(SELF, "IgniteMessage")
			return
		end

		send(SELF,"HandleInteractiveMessage",messagereceived,nil)
	>>
     
    receive InteractiveMessageWithAssignment( string messagereceived, gameSimAssignmentHandle assignment )
     <<
		send(SELF,"HandleInteractiveMessage",messagereceived,assignment)
     >>

	receive HandleInteractiveMessage(string messagereceived, gameSimAssignmentHandle assignment)
	<<
		printl("ai_agent", state.AI.name .. " receive HandleInteractiveMessage: " .. messagereceived)
		
		local setCancelInteraction = false
		
		if messagereceived == "Bury Corpse (player order)" and
			not state.assignment and
			SELF.tags.dead then
			
			send("gameBlackboard",
				"gameObjectRemoveTargetingJobs",
				SELF,
				nil)
			
			send("rendInteractiveObjectClassHandler",
				"odinRendererClearInteractions",
				state.renderHandle)
			
			send("rendCommandManager",
				"odinRendererCreateParticleSystemMessage",
				"Small Beacon",
				state.AI.position.x,
				state.AI.position.y)
			
			if not assignment then
				assignment = query("gameBlackboard",
								"gameObjectNewAssignmentMessage",
								SELF,
								"Burial",
								"",
								"")[1]
			end
			
			send("rendOdinCharacterClassHandler",
				"odinRendererCharacterExpression",
				state.renderHandle,
				"dialogue",
				"jobshovel",
				true)
			
			send( "gameBlackboard",
				"gameObjectNewJobToAssignment",
				assignment,
				SELF,
				"Bury Corpse (player order)",
				"body",
				true )
			
			setCancelInteraction = true
			state.assignment = assignment
			
		elseif messagereceived == "Dump Corpse (player order)" and
			not state.assignment and
			SELF.tags.dead then
			
			SELF.tags.attempt_autoburial = nil
			
			send("gameBlackboard",
				"gameObjectRemoveTargetingJobs",
				SELF,
				nil)
			
			send("rendInteractiveObjectClassHandler",
				"odinRendererClearInteractions",
				state.renderHandle)

			send("rendCommandManager",
				"odinRendererCreateParticleSystemMessage",
				"Small Beacon",
				state.AI.position.x,
				state.AI.position.y)
			
			if not assignment then
				assignment = query("gameBlackboard",
								"gameObjectNewAssignmentMessage",
								SELF,
								"Dump Corpse",
								"",
								"")[1]
			end
			
			send("gameBlackboard",
				"gameObjectNewJobToAssignment",
				assignment,
				SELF,
				"Dump Corpse (player order)",
				"body",
				true )
			
			send("rendOdinCharacterClassHandler",
				"odinRendererCharacterExpression",
				state.renderHandle,
				"dialogue",
				"jobhand",
				true)
			
			setCancelInteraction = true
			state.assignment = assignment
			
		elseif messagereceived == "Cancel corpse orders" and
			SELF.tags.dead then
			send("gameBlackboard",
				"gameObjectRemoveTargetingJobs",
				SELF,
				nil)
			
			SELF.tags.attempt_autoburial = nil
			
			state.assignment = nil
			send(SELF,"resetInteractions")
		end
		
		if setCancelInteraction then
			send("rendInteractiveObjectClassHandler",
                    "odinRendererAddInteractions",
				state.renderHandle,
                         "Cancel orders for corpse of " .. state.AI.name,
                         "Cancel corpse orders",
                         "", --"Cancel corpse orders",
                         "", -- "Cancel corpse orders",
						"graveyard",
						"",
						"Dirt",
						false,true)
		end
	>>
	
	receive ToMidnight()
	<<
		if state.AI.bools.ate_today == false then
			-- do hunger update which assumes no food was eaten today
			-- this is run when someone eats food; you don't want to run it 2x per day
			send(SELF,"updateHungerQoL")
		else
			state.AI.bools.ate_today = false
		end
	
		if not SELF.tags.dead then
			
			if SELF.tags.cult_seed then
				makeMemory("Eldritch: Heard Whispers",nil,nil,nil,nil)
				
			elseif SELF.tags.eldritch_whispers then
				
				makeMemory("Eldritch: Heard Whispers",nil,nil,nil,nil)
				
			elseif SELF.tags.occult_mark_black_magic then
				if rand(1,2) == 1 then
					makeMemory("Eldritch: Rune Visions",nil,nil,nil,nil)
				else
					makeMemory("Eldritch: Inhuman Sigils",nil,nil,nil,nil)
				end
				
			elseif SELF.tags.occult_mark_obeliskian then
				makeMemory("Eldritch: Fell Logic",nil,nil,nil,nil)
				
			elseif SELF.tags.occult_mark_fishperson then
				makeMemory("Eldritch: Beckoning Of The Sea",nil,nil,nil,nil)
				
			elseif SELF.tags.cultist_selenian then
				makeMemory("cultist_selenian",nil,nil,nil,nil)
			elseif SELF.tags.cultist_nature then
				makeMemory("Eldritch: Whispering Trees",nil,nil,nil,nil)
			elseif SELF.tags.cultist_the_queen then
				makeMemory("Eldritch: The Fell Queen",nil,nil,nil,nil)
			elseif SELF.tags.cultist_death then
				makeMemory("Eldritch: Gate of Death",nil,nil,nil,nil)
			end
			
			if SELF.tags.cultist_obeliskian then
				makeMemory("Eldritch: RECEIVED INSTRUCTIONS",nil,nil,nil,nil)
			end
		end
	>>
	
	receive resetInteractions()
	<<
		printl("ai_agent", state.AI.name .. " received resetInteractions")
		
		send("rendInteractiveObjectClassHandler",
				"odinRendererClearInteractions",
				state.renderHandle)
		
		if SELF.tags.dead and
			not SELF.tags.buried and
			not state.assignment then
			
			send("rendInteractiveObjectClassHandler",
                    "odinRendererAddInteractions",
				state.renderHandle,
                         "Give " .. state.AI.name .. " a Proper Burial",
                         "Bury Corpse (player order)",
                         "", --"Bury Corpses",
                         "", --"Bury Corpse (player order)",
						"graveyard",
						"",
						"Dirt",
						false,true)
			
			send("rendInteractiveObjectClassHandler",
                    "odinRendererAddInteractions",
				state.renderHandle,
                         "Dump the Corpse of " .. state.AI.name,
                         "Dump Corpse (player order)",
                         "", -- "Dump Corpses",
                         "", --"Dump Corpse (player order)",
						"graveyard",
						"",
						"Dirt",
						false,true)
			
		end
	>>

	respond getEffectiveSkillLevel(string skillName)
	<<
		-- so we don't have to bake this into every single FSM that queries skill level.
		local effectiveSkill = 1
		
		local valid = false
		for k,v in pairs( EntityDB.HumanStats.skillNameList) do
			if v == skillName then
				valid = true
				break
			end
		end
		
		if valid then
			if state.AI.strs["socialClass"] == "lower" and
				state.AI.currentWorkParty then 
				 -- LC skill inheritance
				 
				local overseer = query( "gameBlackboard",
								"gameObjectGetOverseerMessage",
								state.AI.currentWorkParty)[1]
				
				local overseerAI = query(overseer, "getAIAttributes")[1]
				if overseerAI.skills[skillName] and
					overseerAI.skills[skillName] > 1 then
					
					effectiveSkill = overseerAI.skills[skillName]
				end
				
			else
				if state.AI.skills[skillName] and
					state.AI.skills[skillName] > 1 then
					
					effectiveSkill = state.AI.skills[skillName]
				end
			end
			return "effectiveSkillLevelMessage", effectiveSkill
		else
			printl("ai_agent", "WARNING " .. state.AI.name .. " tried to query unused skill, returning 1" )
			return "effectiveSkillLevelMessage", 1
		end
	>>
	
	respond getAssignment()
	<<
		return "assignmentMessage", state.AI.assignment
	>>
	
	respond getWorkCrew()
	<<
		return "getWorkCrewMessage", state.AI.currentWorkParty
	>>
	
	receive addSelenianHat( gameObjectHandle selenian)
	<<
		-- NOTE: this stuff does NOT work. Well, the hat part anyway.
	
		local models = state.models
		local hatmodel = "models/hats/selenianHeadBlob.upm"

		--[[send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterGeometry", 
			state.renderHandle,
			models["torsoModel"], 
			models["headModel"],
			hatmodel, 
			models["animationSet"],
			"scared")]]
		
		-- TODO push job/animation to flip out and clutch head.
		SELF.tags.selenian_infested = true
		SELF.tags.selenian_infested_stage1 = true
		
		state.AI.traits["Selenian Infested"] = true
		
		-- place the selenian on your head!
		
		local resultROH = query( selenian, "ROHQueryRequest" )[1]
		--[[local model = "models/hats/selenianHeadBlob.upm"
		state.curCarriedHat =selenian
		state.isWearingHat = true

		-- de-register from spatial dictionary
		send("gameSpatialDictionary", 
			 "gridRemoveObject", 
			 selenian )

		-- Start the animation
		send( "rendOdinCharacterClassHandler", 
			  "odinRendererCharacterPickupHatMessage",
			   state.renderHandle,
			   resultROH,
			   "R_H",
			   model)]]
		
	>>
	
	receive transformSelenianHat()
	<<
		-- upgrade to final form
		local models = state.models
		local hatmodel = "models/hats/selenianHeadStalk.upm"
		
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterGeometry", 
			state.renderHandle,
			models["torsoModel"], 
			models["headModel"],
			hatmodel, 
			models["animationSet"],
			"")
	>>
	
	receive removeSelenianHat()
	<<
		local models = state.models
		send("rendOdinCharacterClassHandler",
			"odinRendererSetCharacterGeometry", 
			state.renderHandle,
			models["torsoModel"], 
			models["headModel"],
			models["hatModel"], 
			models["animationSet"],
			"")
		
		state.AI.traits["Selenian Infested"] = false
	>>

	receive updateQoLElement(string elementName,
						int happiness,
						int despair,
						int anger,
						int fear,
						string name,
						string description,
						string icon,
						string iconSkin)
	<<
		state.AI.ints["QoL" .. elementName .. "Happiness"] = happiness
		state.AI.ints["QoL" .. elementName .. "Despair"] = despair
		state.AI.ints["QoL" .. elementName .. "Anger"] = anger
		state.AI.ints["QoL" .. elementName .. "Fear"] = fear
		state.AI.strs["QoL" .. elementName .. "Name"] = name
		state.AI.strs["QoL" .. elementName .. "Description"] = description
		state.AI.strs["QoL" .. elementName .. "Icon"] = icon
		state.AI.strs["QoL" .. elementName .. "IconSkin"] = iconSkin

		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoL" .. elementName .. "Happiness", state.AI.ints["QoL" .. elementName .. "Happiness"])
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoL" .. elementName .. "Despair", state.AI.ints["QoL" .. elementName .. "Despair"])
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoL" .. elementName .. "Anger", state.AI.ints["QoL" .. elementName .. "Anger"])
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoL" .. elementName .. "Fear", state.AI.ints["QoL" .. elementName .. "Fear"])
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoL" .. elementName .. "Name", state.AI.strs["QoL" .. elementName .. "Name"])
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoL" .. elementName .. "Description", state.AI.strs["QoL" .. elementName .. "Description"])
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoL" .. elementName .. "Icon", state.AI.strs["QoL" .. elementName .. "Icon"])
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoL" .. elementName .. "IconSkin", state.AI.strs["QoL" .. elementName .. "IconSkin"])
	>>
	
	receive despawn() override
	<<
		printl("ai_agent", state.AI.name .. "received despawn")
		
		FSM.abort( state, "Despawning.")

		if SELF.tags.fleeing_with_idol then
			printl("ai_agent", "attempting to destroy idol")
			if state.AI.possessedObjects["curPickedUpItem"] then
				local idol = state.AI.possessedObjects["curPickedUpItem"]
				send(idol,"DestroySelf",nil)
			end
			
			printl("ai_agent", "despawning fleeing_with_idol character: " .. state.AI.name )
			send(SELF,"AICancelJob", "despawning")
			send(SELF,"ForceDropEverything")
			send("gameSession", "incSessionInt", "tempCharacterPopulation", -1)
			
			-- No, don't get precious with this. Just delete self and do cleanup.
			-- doing death cleanup.
			sleep()
			
			if state.AI.claimedWorkBuilding then
				send(state.AI.claimedWorkBuilding, "ownerDied", "left")
					local overseer = query("gameBlackboard",
							"gameObjectGetOverseerMessage",
							state.AI.currentWorkParty)[1]
					send("gameBlackboard", "gameSetWorkPartyClaimedBuilding", overseer, nullHandle)
				
				send(state.AI.claimedWorkBuilding,
					"setBuildingOwner",
					nullHandle)
			end
			
			if SELF.tags["cultist"] then
				send("gameSession","incSessionInt", "cultPower", -1)
				if state.group then
					send(state.group,"removeMember", SELF, "left_map")
					if SELF.tags["cult_leader"] then
						send(state.group, "leaderDied", "left the colony")
					end
				end
			end
			
			-- hmm. Not quite appropriate.
			send("gameBlackboard", "gameAgentHasDiedMessage", state.AI, SELF, false)
			send("rendUIManager", "uiRemoveColonist", SELF.id) -- this is good.
			
			state.AI.currentWorkParty = nil
			send("gameSession","incSessionInt","colonyPopulation", -1)
			
			state.AI.bools["canBeSocial"] = false
			send("rendOdinCharacterClassHandler",
				"odinRendererHideCharacterMessage",
				state.renderHandle,
				true)
	
			printl("ai_agent", state.AI.name .. " being removed from game.")
			
			send("rendOdinCharacterClassHandler",
				"odinRendererDeleteCharacterMessage",
				state.renderHandle)
			
			send("gameSpatialDictionary",
				"gridRemoveObject",
				SELF)
			
			send("gameBlackboard",
				"gameObjectRemoveTargetingJobs",
				SELF,
				nil)
			
			destroyfromjob(SELF, nil)
		end
	>>
	
	receive updateSleepQuality( gameObjectHandle building, gameObjectHandle bed)
	<<
		--printl("DAVID", state.AI.name .. " got updateSleepQuality w/ : " .. tostring(building) .. " - " .. tostring(bed) )
		-- aka updateSleepQoL
		
		local rating = 1
		local shortFix = ""
		local longFix = ""
		local description = ""
		local memory = "Slept Outside"
		local effects = {1,0,0,0}
		local ground = false
		local floor = false
		local bed = false
		local building_quality = 0
		
		local icon = "bed"
		local iconSkin = "ui/thoughtIcons.xml"
		
		if building then 
			if state.AI.strs["sleepLocation"] == "ground" then
				ground = true
				icon = "sleep_outside_angry"
				memory = "Slept Outside"
			elseif state.AI.strs["sleepLocation"] == "building" then
				floor = true
				icon = "sleep_on_floor"
				
				if state.AI.strs["socialClass"] == "lower" then	
					memory = "Slept On The Floor"
				elseif state.AI.strs["socialClass"] == "middle" then
					memory = "Slept On The Floor (overseer)"
				elseif state.AI.strs["socialClass"] == "upper" then
					memory = "Slept On The Floor (aristocrat)"
				end
				
			elseif state.AI.strs["sleepLocation"] == "bed" then
				bed = true
				icon = "bed"
				
				if state.AI.strs["socialClass"] == "lower" then
					icon = "lc_bed_icon"
					iconSkin = "ui/orderIcons.xml"
				elseif state.AI.strs["socialClass"] == "middle" then
					icon = "mc_bed_icon"
					iconSkin = "ui/orderIcons.xml"
				elseif state.AI.strs["socialClass"] == "upper" then
					icon = "uc_bed_icon"
					iconSkin = "ui/orderIcons.xml"
				end
				
				building_quality = query(building,"getBuildingQuality")[1]
			end
		end
			
		-- missed sleep entirely for some reason, uh oh.
		local no_sleep = false
		if state.AI.ints.tiredness > 0 then
			no_sleep = true
		
			if state.AI.ints.tiredness >= 1 then
				memory = "A Day Without Sleep"
				icon = "need_sleep"
				iconSkin = "ui/thoughtIcons.xml"
				rating = 1
				shortFix = state.AI.strs.CAPITALIZED_SUBJECTIVE .." needs to sleep!"
				description = state.AI.strs.CAPITALIZED_SUBJECTIVE .." needs to sleep!"
				effects[1] = -12
			elseif state.AI.ints.tiredness >= 2 then
				memory = "Three Days Without Sleep"
				icon = "need_sleep"
				iconSkin = "ui/thoughtIcons.xml"
				rating = 1
				shortFix = state.AI.strs.CAPITALIZED_SUBJECTIVE .." must sleep!"
				description = state.AI.strs.CAPITALIZED_SUBJECTIVE .." must sleep!"
				effects[1] = -18
			elseif state.AI.ints.tiredness >= 3 then
				memory = "Three Days Without Sleep"
				icon = "need_sleep"
				iconSkin = "ui/thoughtIcons.xml"
				rating = 1
				shortFix = state.AI.strs.CAPITALIZED_SUBJECTIVE .." despairs and rages for sleep!"
				effects[1] = -24
				description = state.AI.strs.CAPITALIZED_SUBJECTIVE .." despairs and rages for sleep!"
			end
		end
		
		if state.AI.traits["Pioneering Spirit"] then
			
			if ground then
				rating = 3
				shortFix = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " would be happier sleeping indoors."
				description = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " slept on the ground, which is fine for a Pioneering Spirit, but would be even happier sleeping indoors."
				memory = "Slept Outside (Pioneer)"
			elseif floor then
				shortFix = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " would be happier sleeping in a bed."
				description = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " slept on the floor - much warmer than outside in the dirt! A bed would be even nicer, however."
				memory = "Slept Well"
				rating = 4
			elseif bed then
				rating = 5
				description = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " slept in a bed, how wonderful! Truly this is a luxury for a Pioneering Spirit."
				memory = "Slept Lavishly"
				
			elseif state.AI.ints.tiredness >= 2 then
				rating = 1
				shortFix = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " badly needs an uninterrupted night's sleep!"
				description = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " is very unhappy due to lack of sleep of any kind for days."
				
			elseif state.AI.ints.tiredness >= 1 then
				rating = 2
				shortFix = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " needs an uninterrupted night of sleep."
				description = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " is unhappy due to lack of sleep of any kind, but as a Pioneering Spirit accepts that this sort of thing is part of the job."
			end
			
			-- pioneering spirit is more forgiving.
			local effect_table = { -10, -5, 1, 5, 11}
			effects[1] = effect_table[rating]
			
		elseif state.AI.traits["Materialistic"] and not no_sleep then
			
			if ground then
				rating = 1
				shortFix = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " loathes sleeping on the ground."
				description = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " is extremely upset to be sleeping on the cold, hard ground like some kind of Bandit."
			elseif floor then
				rating = 1
				shortFix = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " loathes sleeping on the floor."
				description = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " is extremely upset to be sleeping on the floor - it might as well be outside in the dirt."
			elseif bed then
				if building_quality < -2 then
					rating = 2
					shortFix = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " wants to sleep in a higher-quality house."
					description = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " slept in a bed, at least. But this home is of appalling quality and is thus upsetting for a Materialistic colonist to sleep in."
					memory = "Slept"
				elseif building_quality >= -2 and building_quality <= 1 then
					rating = 3
					shortFix = state.AI.strs.CAPITALIZED_SUBJECTIVE .." wants to sleep in a higher-quality house."
					description = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " slept in a bed. But this home is of low quality. A Materialistic colonist prefers better-decorated surroundings."
					memory = "Slept Well"
				elseif building_quality > 1 and building_quality < 4 then
					rating = 4
					shortFix = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " wants to sleep in an even higher-quality house."
					description = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " slept in a bed in a reasonably decorated home. This was acceptable, but could be improved - with more decor."
					memory = "Slept Comfortably"
				else
					rating = 5
					description = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " slept in a nice bed in a high quality house filled with material goods; a most pleasing situation."
					memory = "Slept Lavishly"
				end
			end
			
			-- materialist has extreme reactions
			local effect_table = { -18, -7, 0, 4, 14}
			effects[1] = effect_table[rating]
			
		elseif not no_sleep then
			if ground then
				rating = 1
				shortFix = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " needs a building to sleep in - any building."
				description = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " is very upset due to sleeping outside on the ground."
				memory = "Slept Outside"
			elseif floor then
				rating = 2
				
				if state.AI.strs["socialClass"] == "lower" then
					shortFix = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " needs a Cot!"
					description  = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " is upset due to sleeping on a floor."
					memory = "Slept On The Floor"
				elseif state.AI.strs["socialClass"] == "middle" then
					shortFix = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " needs a Practical Bed!"
					description = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " is very upset due to sleeping on a floor like a common Labourer."
					memory = "Slept On The Floor (overseer)"
				elseif state.AI.strs["socialClass"] == "upper" then
					rating = 1
					shortFix = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " needs an Ornate Bed!"
					description = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " is extremely upset due to sleeping on the floor like an animal."
					memory = "Slept On The Floor (aristocrat)"
				end
				
			elseif bed then
				if building_quality < -1 then
					rating = 3
					shortFix = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " wants to sleep in a house with decor."
					description = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " slept in a nice bed, but would be happier if the house had some decor."
					memory = "Slept Well"
				elseif building_quality >= -1 and building_quality <= 2 then
					rating = 4
					shortFix = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " wants to sleep in a house with more decor."
					description = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " slept in a bed in a decent house, but would prefer better decor."
					memory = "Slept Comfortably"
				elseif building_quality > 2 then
					rating = 5
					description = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " slept in a bed in a lovely house."
					memory = "Slept Lavishly"
				end
			end
			
			-- base colonist effects of sleep.
			local effect_table = { -12, -5, 2, 5, 9}
			effects[1] = effect_table[rating]
		end
		
		makeMemory(memory,nil,nil,nil,nil)

		send(SELF,"updateQoLElement","Sleep",
			effects[1], -- happiness, despair, anger, fear
			effects[2],
			effects[3],
			effects[4], 
			memory,
			description,
			icon,
			iconSkin)

		state.AI.ints.QoLSleepRating = rating
		state.AI.strs.QoLSleepHelp = shortFix
		--state.AI.strs.QoLSleepHelpLong = shortFix -- longFix

		send("rendOdinCharacterClassHandler","odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLSleepRating", state.AI.ints.QoLSleepRating)
		send("rendOdinCharacterClassHandler","odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLSleepHelp", state.AI.strs.QoLSleepHelp)
	>>
	
	respond getClaimedWorkBuilding()
	<<
		if state.AI.claimedWorkBuilding then
			return "getClaimedWorkBuildingMessage", state.AI.claimedWorkBuilding
		end
		
		return "getClaimedWorkBuildingMessage", nil
	>>
	
	receive updateSafetyQoL()
	<<
		--[[
			give bonus for first 5 days, starting w/ 5 and counting down.
			default
			no troops, injured
			no troops, uninjured
			barracks exists, at least 1 military
			scaled military
			scaled military 
			"Cowardly" and "Doomed": bump rating scale down one
			"Foolishly Brave": bump rating scale up one
			"Spiritual": require chapel for highest level?
			
			
			-- Proportion of population that is military
			-- Not being injured
			-- Not being inCombat recently
			-- Seeing military patrols recently
			-- Not seeing dead bodies lately
		]]
		
		local emotionValues = {0,0,0,0}  -- happiness, despair, anger, fear
		local icon = "morale4"
		local iconSkin = "ui\\thoughtIcons.xml"
		local description = ""
		local heading = "Safety"
		
		local rating = 1
		local shortFix = ""
		local description = ""
		
		-- military value
		local pop = query("gameSession","getSessionInt","colonyPopulation")[1]
		local milpop = 0
		local results = query("gameSpatialDictionary", "allCharactersWithTagRequest", "military")
		if results and results[1] then
			for k,v in pairs(results[1]) do
				milpop = milpop + 1
			end
		end
		
		local civ_mil_ratio = 0
		if milpop > 0 then
			civ_mil_ratio = math.floor( (milpop / pop) * 100 )
		end
		
		if milpop == 0 then
			shortFix = "Build a Barracks and conscript soldiers!"
			description = description .. state.AI.strs.CAPITALIZED_SUBJECTIVE .. " is fearful and angry due to lack of any military!"
		elseif civ_mil_ratio < 8 then
			shortFix = "Conscript more soldiers to ensure minimal Safety."
			description = description .. state.AI.strs.CAPITALIZED_SUBJECTIVE .. " feels that there could be more soldiers to protect the colony."
			rating = rating + 1
		elseif civ_mil_ratio < 16 then
			shortFix = "Conscript more soldiers to ensure Safety."
			description = description .. state.AI.strs.CAPITALIZED_SUBJECTIVE .. " feels protected by the colonial military."
			rating = rating + 2
		elseif civ_mil_ratio < 24 then
			shortFix = "Conscript more soldiers to ensure complete Safety."
			description = description .. state.AI.strs.CAPITALIZED_SUBJECTIVE .. " feels well-protected by the large colonial military."
			rating = rating + 3
		elseif civ_mil_ratio < 32 then
			description = description .. state.AI.strs.CAPITALIZED_SUBJECTIVE .. " feels extremely well-protected by the huge colonial military."
			rating = rating + 4
		else
			description = description .. state.AI.strs.CAPITALIZED_SUBJECTIVE .. " feels extraordinarily safe due to the enormous military pressence."
			rating = rating + 5
		end
		
		
		local sk_pop = query("gameSession","getSessionInt","steamKnightsActive")[1]
		if sk_pop > 1 then
			rating = rating + 1
			description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE .. " is emboldened by the pressence of a mighty Steam Knight."
		end
		
		local day = query("gameSession","getSessionInt","dayCount")[1]
		if day < 6 then
			local bonus = 6 - day
			rating = rating + bonus
			
			if day < 3 then
				description = "The Colony was recently founded, so the danger of the Frontier hasn't really sunk in yet."
				shortFix = "Satisfied due to recent founding of colony."
			else
				description = "The Colony was recently founded, but the reality of the dangers of the Frontier are starting to sink in."
			end
		end
		
		if state.AI.ints.numAfflictions > 0 then
			shortFix = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " needs healing from a Barber!"
			description = description .. state.AI.strs.CAPITALIZED_SUBJECTIVE .. " is injured and needs healing!"
			rating = rating - 1
		end
		
		if day > 6 then
			-- day > 6, pile it on!
			local isCold = query("gameSession", "getSessionBool", "biomeCold")[1]
			local isDesert = query("gameSession", "getSessionBool", "biomeDesert")[1]
			local isTropical = query("gameSession", "getSessionBool", "biomeTropical")[1]
			
			if isTropical and
				not state.AI.traits["Foolishly Brave"] and
				not state.AI.traits["Pioneering Spirit"] then
				
				rating = rating - 1
				if shortFix == "" then
					shortFix = "Feels endangered by vicious wildlife."
				end
				
				description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE ..
					" finds the constant attacks by tropical wildlife distressing."
					
			elseif isDesert and
				not state.AI.traits["Foolishly Brave"] and
				not state.AI.traits["Hale and Hearty"] and
				not state.AI.traits["Pioneering Spirit"] then
				
				rating = rating - 1
				if shortFix == "" then
					shortFix = "Feels endangered by extreme climate (desert)."
				end
				
				description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE ..
					" fears death by exposure or starvation or bandits and whatever is out there. " .. state.AI.strs.CAPITALIZED_SUBJECTIVE .. " hates the desert."
					
			elseif isCold and
				not state.AI.traits["Foolishly Brave"] and
				not state.AI.traits["Hale and Hearty"] and
				not state.AI.traits["Pioneering Spirit"] then
				
				rating = rating - 1
				if shortFix == "" then
					shortFix = "Feels endangered by extreme climate (cold)."
				end
				
				description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE ..
					" fears death by exposure or starvation or bandits and whatever else plagues this cursed, frozen wasteland."
			end
		end
		
		if rating < 1 then rating = 1 end
		if rating > 5 then rating = 5 end
		
		-- happiness, despair, anger, fear
		local emotionValues = {
			{-12,10,40,20},
			{-4,5,20,10},
			{2,0,0,5},
			{6,0,-5,0},
			{10,0,-10,0}, }
	
		if state.AI.traits["Doomed"] then
			description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE .. " feels Doomed and is thus more fearful regardless of the situation."

			emotionValues = {
				{-16,16,40,30},
				{ -8, 8,20,15},
				{  2, 0, 0,5},
				{  6, 0,-5,2},
				{ 10, 0,-10,0}, }
			 
		elseif state.AI.traits["Craven"] then
			description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE .. " is simply more fearful, whatever the situation, due to being Craven."
			emotionValues = {
				{-12,6,40,40},
				{-4,3,20,20},
				{2,0,0,10},
				{6,0,-5,5},
				{10,0,-10,0}, }
			
		elseif state.AI.traits["Foolishly Brave"] then
			description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE .. " is Foolishly Brave and therefore less upset about any lack of Safety."
			emotionValues = {
				{-6,5,20,10},
				{-2,2,10,5},
				{2,0,0,0},
				{6,0,-5,-5},
				{10,0,-10,-10}, }
			
		end
		
		if rating == 5 then
			description = description .. " "..state.AI.strs.CAPITALIZED_SUBJECTIVE .. " feels completely safe."
			heading = "Feels Completely Safe"
			icon="morale5"
		elseif rating == 4 then
			description = description .. " "..state.AI.strs.CAPITALIZED_SUBJECTIVE .. " feels fairly safe, overall."
			heading = "Colony Is Fairly Safe"
			icon="morale4"
		elseif rating == 3 then
			description = description .. " "..state.AI.strs.CAPITALIZED_SUBJECTIVE .. " feels like safety could be improved.."
			heading = "Could Be More Safe"
			icon="morale3"
		elseif rating == 2 then
			description = description .. " "..state.AI.strs.CAPITALIZED_SUBJECTIVE .. " is feels that Colonial Safety is being neglected."
			heading = "Feels Like Safety Is Ignored"
			icon="morale2"
		elseif rating == 1 then
			description = description .. " "..state.AI.strs.CAPITALIZED_SUBJECTIVE .. " is angry that the Colonial Bureaucrat is so reckless when it comes to Safety."
			heading = "Feels Recklessly Endangered"
			icon="morale1"
		end

		send(SELF,"updateQoLElement",
			"Safety",
			emotionValues[rating][1], -- happiness, despair, anger, fear
			emotionValues[rating][2],
			emotionValues[rating][3],
			emotionValues[rating][4], 
			heading,
			description,
			icon,
			iconSkin)
		
		state.AI.ints.QoLSafetyRating = rating
		state.AI.strs.QoLSafetyHelp = shortFix

		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLSafetyRating", state.AI.ints.QoLSafetyRating)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLSafetyHelp", state.AI.strs.QoLSafetyHelp)
	>>
	
	receive updateWorkQoL()
	<<
		-- do this at noon.
		
		local rating = 1
		local icon = "work"
		local iconSkin = "ui\\thoughtIcons.xml"
		local description = ""
		local shortFix = ""
		local heading = "Work Conditions"
		local overseer = false
		
		--[[	Work quality properties:
			work building quality: -2 to +3
			friend/rival of overseer: -1 to +1
			mood of overseer/workcrew -1 to +1
			various traits: +/- per
			skill: -1 to +2
		]]
		
		if state.AI.strs["socialClass"] == "lower" and state.AI.currentWorkParty then
			overseer = query("gameBlackboard","gameObjectGetOverseerMessage",state.AI.currentWorkParty)[1]
		elseif state.AI.strs["socialClass"] == "lower" and not state.AI.currentWorkParty then
			overseer = false
		else
			overseer = SELF
		end
		
		-- Do feelings about boss.
		-- and mood of boss
		if state.AI.strs["socialClass"] == "lower" and overseer then 
			-- LC and assigned workcrew
			
			-- First, how do we feel about the boss?
			local feelings = query(overseer, "getFeelingsAbout", SELF)[1]
			if not feelings then feelings = 0 end
			
			if feelings > 0 then
				rating = rating + 1
				description = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " likes " .. state.AI.strs.POSSESSIVE .. " boss."
			elseif feelings < 0 then
				rating = rating - 1
				description = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " dislikes " .. state.AI.strs.POSSESSIVE .. " boss."
			end
			
			local overseerAI = query(overseer,"getAIBlock")[1]

			if overseerAI.strs.mood == "happy" then
				rating = rating + 1
				description = description .. " " .. state.AI.strs.CAPITALIZED_POSSESSIVE .. " overseer is happy so work is more pleasant."
			elseif overseerAI.strs.mood == "despair" then
				rating = rating - 1
				description = description .. " " .. state.AI.strs.CAPITALIZED_POSSESSIVE .. " overseer is acting despairingly, which makes work unpleasant."
			elseif overseerAI.strs.mood == "angry" then
				rating = rating - 1
				description = description .. " " .. state.AI.strs.CAPITALIZED_POSSESSIVE .. " overseer is angry, which makes work unpleasant."
			end
			
		elseif state.AI.strs["socialClass"] == "lower" and not overseer then 
			-- LC and no assigned workcrew
			
			if state.AI.traits["Lazy"] or
				state.AI.traits["Of Criminal Element"] then
				
				heading = "Jobless And Loving It"
				rating = rating + 3
				description = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " rather enjoys having no job and would be pleased to loaf about or cause trouble all day."
				
			elseif state.AI.traits["Staunch Traditionalist"] or
				state.AI.traits["Obsequious Bootlicker"] or
				state.AI.traits["Industrious"] then
					
				heading = "Jobless And Upset"
				rating = rating - 1
				description = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " is deeply unhappy about having no job, it being entirely against " .. state.AI.strs.POSSESSIVE .. " character."
				shortFix = "Assign " .. state.AI.strs.OBJECTIVE .. " to a workcrew."
				longFix = "Assign " .. state.AI.strs.OBJECTIVE .. " to a workcrew."
			else
				heading = "Jobless"
				rating = rating - 1
				description = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " is upset about not having a job."
				shortFix = "Assign " .. state.AI.strs.OBJECTIVE .. " to a workcrew."
				longFix = "Assign " .. state.AI.strs.OBJECTIVE .. " to a workcrew."
			end
		end
		
		-- Do work building.
		local workbuilding = false
		if overseer then
			workbuilding = query(overseer,"getClaimedWorkBuilding")[1]
		end
		
		if workbuilding then
			-- working inside, use building quality
			if state.AI.traits["Oblivious Worker"] then
				heading = "Oblivious To Workplace Quality"
				description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE .. " is completely oblivious to the quality (or lack thereof) of their workplace."
				rating = rating +1 
			else
				local buildingQuality = query(workbuilding, "getBuildingQuality")
				if buildingQuality then
					local quality = buildingQuality[1]
					--local descriptor = buildingQuality[2]
					
					if quality >= 6 then
						heading = "Exceptional Workplace"
						description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE .. " is happy to work in a truly exceptional workplace."
						rating = rating + 3
					elseif quality >= 4 then
						heading = "Quality Workplace"
						description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE .. " likes working in a high-quality workplace."
						rating = rating + 2
						if shortFix == "" then
							shortFix = "Improve the quality of " .. state.AI.strs.POSSESSIVE .. " workplace."
						end
					elseif quality >= 2 then
						heading = "Nice Workplace"
						description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE .. " is contented by working in a comfortably-outfitted workplace."
						rating = rating + 1
						if shortFix == "" then
							shortFix = "Improve the quality of " .. state.AI.strs.POSSESSIVE .. " workplace."
						end
					elseif quality >= -1 then
						heading = "Typical Workplace"
						description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE .. " works in a typically-outfitted workplace."
						if shortFix == "" then
							shortFix = "Improve the quality of " .. state.AI.strs.POSSESSIVE .. " workplace."
						end
					elseif quality >= -3 then
						heading = "Uncomfortable Workplace"
						description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE .. " dislikes the uncomfortable workplace " .. state.AI.strs.SUBJECTIVE .. " must work in."
						rating = rating - 1
						if shortFix == "" then
							shortFix = "Improve the quality of " .. state.AI.strs.POSSESSIVE .. " workplace."
						end
					elseif quality >= -5 then
						heading = "Quite Unpleasant Workplace"
						description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE .. " finds " .. state.AI.strs.POSSESSIVE .. " workplace distasteful and unpleasant."
						rating = rating - 1
						if shortFix == "" then
							shortFix = "Improve the quality of " .. state.AI.strs.POSSESSIVE .. " workplace."
						end
					else
						heading = "Wretched Workplace"
						description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE .. " loathes working in this truly wretched place."
						rating = rating - 2
						shortFix = "Improve the quality of " .. state.AI.strs.POSSESSIVE .. " workplace."
					end
				else
					printl("ai_agent", state.AI.name .. " : updateWorkQoL: Error attempting to get buildingQuality info!")
				end
			end
			
		elseif overseer then
			-- working outdoors.
			icon = "jobtree"
			-- working outside
			if state.AI.traits["Interest in Exotic Wildernesses"] or
				state.AI.traits["Pioneering Spirit"] or
				state.AI.traits["Rustic Disposition"] or
				state.AI.traits["Woodtouch"] or
				state.AI.traits["Hale and Hearty"] then
				
				if heading == "Work Conditions" then
					heading = "Enjoys Working Outdoors"
				end
				
				description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE .. " enjoys the outdoors and is pleased to be working outside."
				rating = rating +3
			else
				
				if query("gameSession","getSessionBool","outdoor_qol_bonus_unlocked")[1] == true then --you have a tech bonus!
					if heading == "Work Conditions" then
						heading = "Enjoys Working Outdoors due to Indoctrination"
					end
					description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE .. " is convinced of the merits of working outdoors by our indoctrination programme."
					rating = rating +2
				else
					if heading == "Work Conditions" then
						heading = "Working Outdoors"
					end
					description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE .. " has no strong opinion about working outside."
				end
			end
		end

		if SELF.tags.military and
			query("gameSession","getSessionBool","military_happy1_unlocked")[1] then
			
			description = description .. " This soldier is happier due to the 'Stiff Upper Lip' discovery."
			rating = rating +1
		end
		
		local isCold = query("gameSession", "getSessionBool", "biomeCold")[1]
		local isDesert = query("gameSession", "getSessionBool", "biomeDesert")[1]
		local isTropical = query("gameSession", "getSessionBool", "biomeTropical")[1]

		if isDesert and
			not state.AI.traits["Hale and Hearty"] and
			not state.AI.traits["Pioneering Spirit"] then
			
			rating = rating - 1
			if shortFix == "" then
				shortFix = "Dislikes work in extreme climate (desert)."
			end
			description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE ..
				" dislikes working in the desert due to the stifling heat."
				
		elseif isTropical and
			not state.AI.traits["Hale and Hearty"] and
			not state.AI.traits["Interest in Exotic Wildernesses"] and
			not state.AI.traits["Pioneering Spirit"] then
			
			rating = rating - 1
			if shortFix == "" then
				shortFix = "Dislikes work in extreme climate (tropics)."
			end
			description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE ..
				" dislikes working in the tropics due to the stifling heat."
				
		elseif isCold and
			not state.AI.traits["Hale and Hearty"] and
			not state.AI.traits["Pioneering Spirit"] then
			
			rating = rating - 1
			if shortFix == "" then
				shortFix = "Suffers from working in extreme climate (cold)."
			end
			
			description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE ..
				" hates working in this climate where it is always bone-chillingly cold."
		end
		
		if rating < 1 then rating = 1 end
		if rating > 5 then rating = 5 end
		
		-- happiness, despair, anger, fear
		local emotionValues = {
			{-20,6,6,0},
			{-10,2,2,0},
			{0,0,0,0},
			{9,0,0,0},
			{18,-4,-4,0},
		}
		
		send(SELF,"updateQoLElement",
				"WorkConditions",
				emotionValues[rating][1], -- happiness, despair, anger, fear
				emotionValues[rating][2],
				emotionValues[rating][3],
				emotionValues[rating][4], 
				heading,
				description,
				icon,
				iconSkin)
		
		state.AI.ints.QoLWorkConditionsRating = rating
		state.AI.strs.QoLWorkConditionsHelp = shortFix
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLWorkConditionsRating", state.AI.ints.QoLWorkConditionsRating)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLWorkConditionsHelp", state.AI.strs.QoLWorkConditionsHelp)
	>>
	
	respond isOnShift()
	<<
		local onShift = query("gameBlackboard","gameObjectGetWorkPartyOnShift",state.AI.currentWorkParty)[1]
		return "isOnShiftMessage", onShift
	>>
	
	receive updateCrowdingQoL()
	<<
		-- very simple for now, crowding starts at all stats at +5 good stuff to -10 bad stuff
		-- depending on where you are in scale of 1 to 100 pop.
		local rating = 5
		local shortFix = ""
		local longFix = ""
		local description = ""
		local heading = ""
		local pop = query("gameSession","getSessionInt","colonyPopulation")[1]
		
		local crowdingDescriptions = {
			"It is positively, insufferably squalorous.",
			"It is positively, insufferably squalorous.",
			"Lots of new riff-raff about, aren't there.",
			"Not terrifically crowded. Almost cozy.",
			"Pleasantly sparse; the air is fresh and land ripe for taking.",
		}
		local crowdingHeadings = {
			"Upsettingly Crowded",
			"Upsettingly Crowded",
			"Acceptably Populated",
			"Decently Uncrowded",
			"Pleasantly Sparse",
		}
		
		if state.AI.traits["Gregarious"] then
			-- loves people.
			crowdingDescriptions = {
				"It's positively lonely and barren. More people should be about!",
				"It's rather lonely and conversation-starved here. Where is everyone?",
				"It's starting to get some life around here, though wouldn't mind some more friends.",
				"It's always a party, lots of people to meet and talk to; lovely!",
				"It's always a party, lots of people to meet and talk to; lovely!",
			}
			crowdingHeadings = {
				"Lonely and Awful",
				"Suffocatingly Sparse",
				"Decently Populated",
				"Delightfully Crowded",
				"Delightfully Crowded",
			}
			rating = 5 - math.floor( pop * 0.3333 )
			if rating < 2 then
				shortFix = state.AI.strs.CAPITALIZED_SUBJECTIVE .. " would be happier if the population was higher."
			end
			
			if pop > 135 then
				rating = 5
			elseif pop > 105 then
				rating = 4
			elseif pop > 75 then
				rating = 3
			elseif pop > 45 then
				rating = 2
			else
				rating = 1
			end
			
			description = crowdingDescriptions [ rating ]
			heading = crowdingHeadings[ rating ]
			
		elseif state.AI.traits["Hermit"] then
			-- Doesn't like people.
			if pop > 70 then
				rating = 1
			elseif pop > 60 then
				rating = 2
			elseif pop > 40 then
				rating = 3
			elseif pop > 25 then
				rating = 4
			else
				rating = 5
			end
			
			description = crowdingDescriptions [ rating ]
			heading = crowdingHeadings[ rating ]
			
		else
			-- normal person.
			if pop > 135 then
				rating = 1
			elseif pop > 105 then
				rating = 2
			elseif pop > 75 then
				rating = 3
			elseif pop > 45 then
				rating = 4
			else
				rating = 5
			end
			
			description = crowdingDescriptions [ rating ]
			heading = crowdingHeadings[ rating ]
			
			--rating = 5 - math.floor( pop * 0.3333 )
		end
		
		if state.AI.traits["Gregarious"] then
			description = description .. " Due to " .. state.AI.strs.POSSESSIVE .. " Gregarious nature, " ..
				state.AI.strs["firstName"] .. " feels more comfortable in a crowded colony."
				
		elseif state.AI.traits["Hermit"] then
			
			description = description .. " Due to " .. state.AI.strs.POSSESSIVE .. " being something of a Hermit, " ..
				state.AI.strs["firstName"] .. " feels more comfortable in a low-population colony."
		end
		
		if rating < 1 then rating = 1 end
		if rating > 5 then rating = 5 end
		
		local isCold = query("gameSession", "getSessionBool", "biomeCold")[1]
		if isCold then
			rating = rating - 1
			if shortFix == "" then
				shortFix = "Disastisfied due to cold climate."
			end
			description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE ..
				" finds that the cold makes people grumpier and enjoying life difficult."
		end
		
		send(SELF,"updateQoLElement",
				"Crowding",
				(rating * 2) - 5, -- happiness, despair, anger, fear
				0,
				(3 - rating)  * 3, -- div( pop, 8),
				0,
				heading,
				description,
				"population_icon",
				"ui/orderIcons.xml")
		
		state.AI.strs.QoLCrowdingHelp = shortFix
		state.AI.strs.QoLCrowdingHelpLong = longFix
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLCrowdingRating", state.AI.ints.QoLCrowdingRating)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLCrowdingHelp", state.AI.strs.QoLCrowdingHelp)
	>>
	
	receive updateHungerQoL()
	<<
		local classfooddefault = 1
		if state.AI.strs.socialClass == "middle" then
			classfooddefault = 2
		elseif state.AI.strs.socialClass == "upper" then
			classfooddefault = 3
		end
		
		local rating = 1
		local shortFix = ""
		local heading = "Hunger Satisfied"
		local icon = "food_okay"
		local iconSkin = "ui/thoughtIcons.xml"
		local description = state.AI.strs.CAPITALIZED_SUBJECTIVE
		
		-- move all values down one, clear current day, collate quality of last 3 days vs. social class
		if state.AI.bools.ate_today == false then
			table.insert(state.foodrecord, 1, 0 ) -- default is starving; LC more used to starving.
		end
		table.remove(state.foodrecord, #state.foodrecord) -- and kill last record.
		
		local went_without_food = false
		local food_points = 0
		for k,v in pairs(state.foodrecord) do
			--printl("DAVID", " doing foodpoints k/v = " .. tostring(k) .. " / " .. tostring(v) )
			if v == 0 then
				-- no food for a day, very bad!
				food_points = food_points -1
			elseif v == classfooddefault then
				food_points = food_points + 1
			elseif v > classfooddefault then
				food_points = food_points + 2
			end
			-- else had food of "bad" quality which gives 0 point change.
		end
		--printl("DAVID", state.AI.name .. " food points before correction = " .. food_points)
		food_points = math.ceil( food_points * 0.3333 )
		--printl("DAVID", state.AI.name .. " food points AFTER correction = " .. food_points)
		
		if SELF.tags.starving then
			description = description .. " is starving! Find food, any food - quickly!"
			shortFix = "Starving! Find food immediately."
			heading = "Tortured By Hunger"
			icon = "cannibalism"
			
		elseif food_points == -1 then
			description = description .. " has gone without food in the last few days! This is very upsetting."
			shortFix = "Gone without food recently! Find some food."
			heading = "Suffering Due To Hunger"
			icon= "hungry"
		elseif food_points == 0 then
			description = description .. " has been eating bad-quality food for the last few days and is unhappy."
			shortFix = "Could use better-quality food."
			heading = "Poorly Fed"
			icon= "food_bad"
		elseif food_points == 1 then
			description = description .. " has been satisfactory fed for the last few days."
			shortFix = "Well-fed but happier with higher-quality food."
			heading = "Well-Fed But Unsated"
			
		elseif food_points == 2 then
			description = description .. " has eaten above " .. state.AI.strs.POSSESSIVE .. " station and feels very well taken care of."
			icon= "food_plate"
			heading = "Feasting Like A King"
		end
		
		rating = rating + food_points
	
		if state.AI.traits["Epicurean"] then
			description = description .. " As an Epicurean, " .. state.AI.strs.OBJECTIVE .. " has higher standards for food and drink than most colonists."
			rating = rating -1
		end
		
		-- move all values down one, clear current day, collate quality of last 3 days vs. social class
		if SELF.tags.drank_today then
			-- then ... good.
			SELF.tags.drank_today = nil
		else
			-- no drink today, insert blank record.
			table.insert(state.drinkrecord, 1, 0 )
			-- and kill last record.
			table.remove(state.drinkrecord, #state.drinkrecord) 
		end
		
		local had_tea = false
		local had_brew = false
		local had_spirits = false
		for k, v in pairs(state.drinkrecord) do
			if v == 1 then
				had_tea = true
			elseif v == 2 then
				had_brew = true
			elseif v == 3 then
				had_spirits = true
			end
		end
		
		if had_spirits then
			description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE ..
				" drank some potent spirits at the Pub recently and can't remember being particularly upset about anything."
			
			rating = rating + 3
		elseif had_brew then
			description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE ..
				" drank brew at the Pub recently, which helped " .. state.AI.strs.OBJECTIVE ..
				" feel better about the difficulties of Frontier Life."
				
			rating = rating + 2
			
			if shortFix == "" and rating < 5 then
				shortFix = "Could use a more potent drink."
			end
			
		elseif had_tea then
			description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE ..
				" had a cup of tea recently, which is Civilized and Proper. "
			
			rating = rating +1
			if shortFix == "" and rating < 5 then
				shortFix = "Could use a proper drink at a Pub."
			end
		else
			-- no pub?
			description = description .. " " .. state.AI.strs.CAPITALIZED_SUBJECTIVE ..
				" is unhappy that the Colony doesn't have a proper Public House to serve drinks. "
				
			if shortFix == "" then
				shortFix = "Wants a Public House."
			end
		end
		
		if rating < 1 then rating = 1 end
		if rating > 5 then rating = 5 end
		
		-- happiness, despair, anger, fear
		local emotionValues = {
				{-24,32,6,0},
				{-12,14,3,0},
				{0,0,0,0},
				{5,0,0,0},
				{10,-5,0,0},
			}
		
		send(SELF,"updateQoLElement",
				"Hunger",
				emotionValues[rating][1], 
				emotionValues[rating][2],
				emotionValues[rating][3],
				emotionValues[rating][4], 
				heading,
				description,
				icon,
				iconSkin)
		
		state.AI.ints.QoLHungerRating = rating
		state.AI.strs.QoLHungerHelp = shortFix
		send("rendOdinCharacterClassHandler", "odinRendererCharacterSetIntAttributeMessage", state.renderHandle, "QoLHungerRating", state.AI.ints.QoLHungerRating)
		send("rendOdinCharacterClassHandler", "odinRendererSetCharacterAttributeMessage", state.renderHandle, "QoLHungerHelp", state.AI.strs.QoLHungerHelp)
	>>
	
	receive resetMeleeWeapon()
	<<
		-- make sure we're using the melee weapon appropriate to our profession
		-- Note: soldiers will use melee attack on ranged weapons if said weapon has a melee attack ability
		local weapon = state.AI.strs.melee_weapon
		-- resets weapon if not using ANY weapon
		-- weapon is properly reset along with profession via changeProfession function
		if not weapon then
			local profession = state.AI.strs["citizenClass"]
			local data = EntityDB[profession]
			if data.melee_weapons then
				-- set random? default weapon
				local random_weapon = data.melee_weapons[ rand(1,#data.melee_weapons)]
				
				--printl("ai_agent", state.AI.name .. " setting melee wep to: " .. random_weapon)
				send(SELF,"setWeapon","melee", random_weapon )
			end
		end
	>>
	
	receive resetRangedWeapon()
	<<
		-- make sure we're using the ranged weapon appropriate to our profession / barracks (if applicable)
		local weapon = state.AI.strs.ranged_weapon
		local profession = state.AI.strs["citizenClass"]
		
		if not weapon then
			local data = EntityDB[profession]
			if data.ranged_weapons then
				-- set random? default weapon
				send(SELF,"setWeapon","ranged",data.ranged_weapons[ rand(1,#data.ranged_weapons)] )
			end
		end
		
		if SELF.tags.military then
			
			-- make sure you're using the correct weapon for your barracks. If you have a barracks.
			-- AND make sure barracks has ammo for weapon. If not, default to pistol.
			if state.AI.strs["socialClass"] == "lower" and state.AI.currentWorkParty then 
				local overseer = query("gameBlackboard",
							    "gameObjectGetOverseerMessage",
							    state.AI.currentWorkParty)[1]
				
				local rax = query(overseer, "getClaimedWorkBuilding")[1]
				if not rax then
					if weapon then
						if weapon ~= "pistol" then
							send(SELF,"setWeapon","ranged","pistol")
						end
					else
						send(SELF,"setWeapon","ranged","pistol")
					end
				else
					-- get loadout weapon from barracks
					local rax_tags = query(rax,"getTags")[1]
					weapon = query(rax,"getWeaponLoadout")
					if weapon and weapon[1] then
						weapon = weapon[1]
						
						local weapon_data = EntityDB[weapon]
						if weapon_data.ammo_tier then
							if rax_tags["no_supplies" .. tostring(weapon_data.ammo_tier)] == true then
								-- no ammo, default to pistol.
								weapon = "pistol"
							end
						end
						if state.AI.strs.ranged_weapon ~= weapon then
							send(SELF,"setWeapon","ranged",weapon)
						end
					end
				end
				
			elseif state.AI.strs["socialClass"] == "middle" then
				if not state.AI.claimedWorkBuilding then
					if weapon then
						if weapon ~= "pistol" then
							send(SELF,"setWeapon","ranged","pistol")
						end
					else
						send(SELF,"setWeapon","ranged","pistol")
					end
				else
					-- get loadout weapon from barracks
					local rax_tags = query(state.AI.claimedWorkBuilding,"getTags")[1]
					weapon = query(state.AI.claimedWorkBuilding,"getWeaponLoadout")
					if weapon and weapon[1] then
						weapon = weapon[1]
						
						local weapon_data = EntityDB[weapon]
						if weapon_data.ammo_tier then
							if rax_tags["no_supplies" .. tostring(weapon_data.ammo_tier)] == true then
								-- no ammo, default to pistol.
								weapon = "pistol"
							end
						end
						if state.AI.strs.ranged_weapon ~= weapon then
							send(SELF,"setWeapon","ranged",weapon)
						end
					end
				end 
			end
			
		end
	>>

	receive incSkill(string skillname, int amount)
	<<
		printl("ai_agent", state.AI.name .. " received incSkill, amount = " .. tostring(amount))
		if SELF.tags["middle_class"] then	
			if state.AI.skills[skillname] and state.AI.skills[skillname] < 5 then
				state.AI.skills[skillname] = state.AI.skills[skillname] + amount
				if state.AI.skills[skillname] >= 5 then
					state.AI.skills[skillname] = 5
				elseif state.AI.skills[skillname] < 0 then
					state.AI.skills[skillname] = 0
				end
				
				send("rendOdinCharacterClassHandler",
					"odinRendererCharacterSetIntAttributeMessage",
					state.renderHandle,
					skillname,
					state.AI.skills[skillname])
			end
		end
		
		if state.AI.claimedWorkBuilding then
			send(state.AI.claimedWorkBuilding, "refreshSkillDisplay" )
		end
	>>
	
	receive refreshMaxHealth()
	<<
		local maxHealthBonus = query("gameSession", "getSessionInt", "civilianHealthBonus")[1]
		
		if state.AI.traits["Hale and Hearty"] then
			maxHealthBonus = maxHealthBonus + 2
		end
		
		if maxHealthBonus ~= 0 then
			local humanstats = EntityDB["HumanStats"]
			state.AI.ints["healthMax"] = humanstats["healthMax"] + maxHealthBonus
		end
	>>
	
	receive rudelyWakeFromSleep()
	<<
		if not state.AI.traits["Heavy Sleeper"] then
			send(SELF,"addTag","rudely_awoken")
			if state.AI.bools["asleep"] then
				-- sleep.fsm will detect when this bool is flipped and do the correct wakeup sequence
				state.AI.bools["asleep"] = false
			end
		end
	>>
	
	receive registerWithHouse( gameObjectHandle house)
	<<
		
	>>
	
	receive refreshCharacterAlert()
	<<
		if SELF.tags.dead then
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "no_overseer", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "starvation", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "frontier_justice", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "rage_state", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "despair_sate", "", "")
			
		elseif SELF.tags.frontier_justice then
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "no_overseer", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "starvation", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "frontier_justice",
				state.AI.name .. " has been ordered to be shot.",
				"frontier_justice")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "rage_state", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "despair_sate", "", "")
			
		elseif SELF.tags.starving then
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "no_overseer", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "starvation",
				state.AI.name .. " is starving and will soon die without food.",
				"cannibalism")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "frontier_justice", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "rage_state", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "despair_sate", "", "")
			
		elseif SELF.tags.rage_state then
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "no_overseer", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "starvation", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "frontier_justice", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "rage_state",
				state.AI.name .. " is in a rage and will not work until calmed down.",
				"fist_upheld_red")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "despair_sate", "", "")
			
		elseif SELF.tags.despair_state then
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "no_overseer", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "starvation", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "frontier_justice", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "rage_state", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "despair_sate",
				state.AI.name .. " is Maddened due to high Despair and is prone to irrational decisions.",
				"madness_extreme")
			
		elseif not state.AI.currentWorkParty and not SELF.tags.upper_class then
			SELF.tags.unassigned_labourer = true
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "no_overseer",
				"",
				"questionmark")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "starvation", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "frontier_justice", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "rage_state", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "despair_sate", "", "")
		else
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "no_overseer", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "starvation", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "frontier_justice", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "rage_state", "", "")
			send("rendOdinCharacterClassHandler", "odinRendererCharacterSetAlert", state.renderHandle, "despair_sate", "", "")
		end
		
		if state.AI.currentWorkParty and not SELF.tags.upper_class then
			SELF.tags.unassigned_labourer = nil
		end
	>>
	
	receive BoardVehicle ( gameObjectHandle gOH )
	<<
		state.AI.thinkLocked = true
		printl("boarding vehicle GOH");
		resultROH = query( gOH, "ROHQueryRequest" );
		printl("boarding vehicle GOH: got a resultROH "..resultROH[1]);
		seatingArrangement = query ( gOH, "seatingRequest" );
		printl("BoardVehicle: Got seat " .. seatingArrangement[1]);
		state.seat = seatingArrangement[1]

		SELF.tags["combat_target_for_enemy"] = false
		-- FIXME: obviously, we will have to set these from the entity in question

		send("gameSpatialDictionary", "gridRemoveObject", SELF);

		send(gOH, "BoardVehicle2", SELF, state.renderHandle, seatingArrangement[1] );

	>>
	
	receive UnboardVehicle ( gameObjectHandle gOH, int offsetX, int offsetY )
	<<
		-- don't like "disembark"?
		state.AI.thinkLocked = false
		newPosition = query(gOH, "gridGetPosition");
		newPosition[1].x = newPosition[1].x + offsetX
		newPosition[1].y = newPosition[1].y + offsetY
		resultROH = query( gOH, "ROHQueryRequest" );
		send("rendOdinCharacterClassHandler", "odinRendererCharacterDropCharacterMessage", resultROH[1], state.renderHandle, state.seat, "", "", "" );

		setposition(newPosition[1].x,newPosition[1].y);
		send("rendOdinCharacterClassHandler", 
				"odinRendererTeleportCharacterMessage", 
				state.renderHandle, 
				newPosition[1].x,
				newPosition[1].y );	
			send("rendOdinCharacterClassHandler",
			"odinRendererIdleCharacterMessage",
			state.renderHandle);

		-- NV FIXME: Trader Goods on a boat should not be necessarily treated this way.

		printl( state.AI.strs["firstName"] .. " " .. state.AI.strs["lastName"]  .. ": gameObject " .. tostring( SELF ) .. " (rOH:" .. state.renderHandle .. ") at " .. tostring( state.AI.position ) )
		SELF.tags["combat_target_for_enemy"] = true
	>>
>>