import json
import xml.etree.ElementTree as ET

class Res():
    def __init__(self, args, cfg, scene_id):
        self.args = args
        self.cfg = cfg
        self.scene_id = scene_id
        
        self.map = None
        self.food = None
        
        self.load()
    
    def load(self):
        # config/terrain/map.json
        try:
            mapfile = "%s/%d.json" % (self.cfg["terrain"], self.scene_id)
            f = open(mapfile, 'rt')
            self.map = json.loads(f.read())
            f.close()
        except Exception as e:
            print(e)
            exit(0)
        
        # config/xml/food.xml
        try:
            foodfile = "%s/food.xml" % (self.cfg["xml"])
            root = ET.parse(foodfile).getroot()
            
            self.food = {}
            for food in root.findall("food"): #<food mapid="1002">
                mapid = food.get("mapid")
                if mapid == str(self.scene_id):
                    for item in food.findall("item"): #<item id="1" type="11" size="0.2" mapnum="400"/>
                        self.food[int(float(item.get("id")))] = { "size": float(item.get("size")), "type": float(item.get("type")) }
                    break
        except Exception as e:
            print(e)
            exit(0)
        #print(self.food)


g_res = None
def new(args, cfg, scene_id):
    global g_res
    g_res = Res(args, cfg, scene_id)
    return g_res