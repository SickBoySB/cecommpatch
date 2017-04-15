printl("SETTING DESERT BIOME")
local biomeName = "desert"
send("gameSession", "setSessionString", "biome", biomeName)
send("gameSession", "setSessionBool", "spawnExtraBandits", true)
send("gameSession", "setSessionBool", "biomeDesert", true)
send("gameSession", "setSessionBool", "biomeCold", false)
send("gameSession", "setSessionBool", "biomeTemperate", false)
send("gameSession", "setSessionBool", "biomeTropical", false)
send("gameSession", "setSessionInt", "colonyPopulation", 0)
send("gameSession", "setSessionInt", "militaryCount", 0)

send("rendCommandManager", "odinRendererSetLighting",118,79,75,38,27,53) -- Desert Lighting

-- set up agriculture for this biome
-- (these will be read by farm.go to set allowed crops)
local cropTable = EntityDB.WorldStats.climateInfoPerBiome[ biomeName ].cropTable
for cropName, stats in pairs( cropTable ) do
     send("gameSession", "setSessionBool", "cropUnlocked=" .. cropName, cropTable[cropName].unlocked )
	send("gameSession", "setSessionInt", "cropGrowthModifier=" .. cropName, cropTable.growthModifier)
end

function spawnGameobject( x, y, objectType, objectTable )
	if x > 235 then x = 235 end
	if x < 20 then x = 20 end
	
	if y > 235 then y = 235 end
	if y < 20 then y = 20 end
	
	local createResults = query("scriptManager",
						   "scriptCreateGameObjectRequest",
						   objectType,
						   objectTable )
	
	local handle = createResults[1]
	if handle ~= nil then
		send(handle, "GameObjectPlace", x, y )
	end
end

-- need some accessible hunting at game start
local animals_to_spawn = { [1] = {["legacyString"]="Beetle"}, [2] = {["legacyString"]="Desert Fox"} }

--spawnGameobject( 250, 210, "herd", animals_to_spawn[rand(1,#animals_to_spawn)])
--spawnGameobject( 250, 300, "herd", animals_to_spawn[rand(1,#animals_to_spawn)])
spawnGameobject( rand(20,235), rand(20,235), "herd", animals_to_spawn[rand(1,#animals_to_spawn)])
spawnGameobject( rand(20,235), rand(20,235), "herd", animals_to_spawn[rand(1,#animals_to_spawn)])

-- load in some more Obeliskians! Fun fun fun.

local x_max = 255
local y_max = 255

spawnGameobject( 	rand(20, 98),
				rand(20,235),
				"objectcluster",
				{ legacyString="Obeliskian Cluster" } )
		
spawnGameobject( 	rand(20,235),
				rand(20, 98),
				"objectcluster",
				{ legacyString="Obeliskian Cluster" } )

