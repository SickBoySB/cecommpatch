printl("SETTING COLD BIOME")
local biomeName = "cold"
send("gameSession", "setSessionString", "biome", biomeName)
send("gameSession", "setSessionBool", "spawnExtraBandits", true)
send("gameSession", "setSessionBool", "biomeCold", true)
send("gameSession", "setSessionBool", "biomeTemperate", false)
send("gameSession", "setSessionBool", "biomeTropical", false)
send("gameSession", "setSessionBool", "biomeDesert", false)
send("gameSession", "setSessionInt", "colonyPopulation", 0)
send("gameSession", "setSessionInt", "militaryCount", 0)

send("rendCommandManager", "odinRendererSetLighting",111,90,59,93,83,69) -- Arctic Lighting

-- set up agriculture for this biome
-- (these will be read by farm.go to set allowed crops)

local cropTable = EntityDB.WorldStats.climateInfoPerBiome[ biomeName ].cropTable
for cropName, stats in pairs( cropTable ) do
     send("gameSession", "setSessionBool", "cropUnlocked=" .. cropName, cropTable[cropName].unlocked )
	send("gameSession", "setSessionInt", "cropGrowthModifier=" .. cropName, cropTable.growthModifier)
end

function spawnGameobject( x, y, objectType, objectTable )
	local createResults = query( "scriptManager", "scriptCreateGameObjectRequest", objectType, objectTable )
	local handle = createResults[1]
	if handle ~= nil then
		send(handle, "GameObjectPlace", x, y )
	end
end

-- need some accessible hunting at game start
local animals_to_spawn = { [1] = {["legacyString"]="Wooly Aurochs"}, [2] = {["legacyString"]="Arctic Dodo"} }
--spawnGameobject( 250, 210, "herd", animals_to_spawn[rand(1,#animals_to_spawn)])
--spawnGameobject( 250, 300, "herd", animals_to_spawn[rand(1,#animals_to_spawn)])
spawnGameobject( rand(20,235), rand(20,235), "herd", animals_to_spawn[rand(1,#animals_to_spawn)])
spawnGameobject( rand(20,235), rand(20,235), "herd", animals_to_spawn[rand(1,#animals_to_spawn)])