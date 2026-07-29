$version: "2"

namespace com.amazonaws.sdk.benchmark

use aws.protocols#awsJson1_0
use smithy.protocols#rpcv2Cbor
use smithy.test#httpRequestTests
use smithy.test#httpResponseTests

@documentation("""
    From Amazon DynamoDB.
    Deserialization of recursive structures.
""")
@http(method: "POST", uri: "/GetItem", code: 200)
@httpRequestTests([
    {
        id: "awsJson1_0_GetItemInput_Baseline"
        protocol: awsJson1_0
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Key: {
                id: {
                    S: "test-id"
                }
            }
        }
        tags: ["serde-benchmark"]
    }
])
@httpResponseTests([
    // section: json
    {
        id: "awsJson1_0_GetItemOutput_Baseline"
        protocol: awsJson1_0
        code: 200
        body: """
        {}
        """
        tags: ["serde-benchmark"]
    }
    {
        id: "awsJson1_0_GetItemOutput_S"
        protocol: awsJson1_0
        code: 200
        body: """
{
  "Item": {
    "id": {
      "S": "recipe-001"
    },
    "name": {
      "S": "Classic Carbonara"
    },
    "cuisine": {
      "S": "Italian"
    },
    "cook_time": {
      "N": "20"
    },
    "difficulty": {
      "S": "Medium"
    },
    "rating": {
      "N": "4.8"
    }
  },
  "ConsumedCapacity": {
    "TableName": "pasta-recipes",
    "CapacityUnits": 1.1
  }
}
"""
        tags: ["serde-benchmark"]
    }
    {
        id: "awsJson1_0_GetItemOutput_M"
        protocol: awsJson1_0
        code: 200
        body: """
{
  "Item": {
    "id": {
      "S": "recipe-002"
    },
    "name": {
      "S": "Fettuccine Alfredo"
    },
    "description": {
      "S": "Creamy, rich pasta dish with butter, parmesan cheese, and fresh fettuccine noodles"
    },
    "cook_time": {
      "N": "25"
    },
    "prep_time": {
      "N": "15"
    },
    "difficulty": {
      "S": "Easy"
    },
    "cuisine": {
      "S": "Italian"
    },
    "servings": {
      "N": "4"
    },
    "rating": {
      "N": "4.6"
    },
    "tags": {
      "SS": ["creamy", "comfort-food", "vegetarian"]
    },
    "ingredients": {
      "L": [
        {
          "M": {
            "item": {
              "S": "fettuccine pasta"
            },
            "amount": {
              "S": "1 lb"
            }
          }
        },
        {
          "M": {
            "item": {
              "S": "butter"
            },
            "amount": {
              "S": "1/2 cup"
            }
          }
        },
        {
          "M": {
            "item": {
              "S": "parmesan cheese"
            },
            "amount": {
              "S": "1 cup grated"
            }
          }
        },
        {
          "M": {
            "item": {
              "S": "heavy cream"
            },
            "amount": {
              "S": "1/2 cup"
            }
          }
        }
      ]
    },
    "nutrition": {
      "M": {
        "calories": {
          "N": "520"
        },
        "protein": {
          "N": "18"
        },
        "carbs": {
          "N": "45"
        },
        "fat": {
          "N": "28"
        }
      }
    }
  },
  "ConsumedCapacity": {
    "TableName": "pasta-recipes",
    "CapacityUnits": 2.5,
    "ReadCapacityUnits": 2.5
  }
}
"""
        tags: ["serde-benchmark"]
    }
    {
        id: "awsJson1_0_GetItemOutput_L"
        protocol: awsJson1_0
        code: 200
        body: """
{
  "Item": {
    "id": {
      "S": "recipe-003"
    },
    "name": {
      "S": "Grandma's Ultimate Lasagna Bolognese"
    },
    "description": {
      "S": "A traditional Italian lasagna recipe passed down through generations, featuring layers of rich meat sauce, creamy bechamel, fresh pasta sheets, and a blend of artisanal cheeses. This complex dish requires multiple preparation stages and represents the pinnacle of Italian comfort food craftsmanship. Recipe adapted from 'La Cucina della Nonna' by Maria Benedetti, 1952."
    },
    "cook_time": {
      "N": "180"
    },
    "prep_time": {
      "N": "120"
    },
    "total_time": {
      "N": "300"
    },
    "difficulty": {
      "S": "Expert"
    },
    "cuisine": {
      "S": "Italian"
    },
    "servings": {
      "N": "12"
    },
    "rating": {
      "N": "4.9"
    },
    "cost_estimate": {
      "N": "45.50"
    },
    "active": {
      "BOOL": true
    },
    "featured": {
      "BOOL": true
    },
    "tags": {
      "SS": ["traditional", "comfort-food", "family-recipe", "holiday", "meat-sauce", "layered", "baked", "italian-classic", "time-intensive", "special-occasion"]
    },
    "categories": {
      "SS": ["main-course", "pasta", "casserole", "italian"]
    },
    "allergens": {
      "SS": ["dairy", "gluten", "eggs"]
    },
    "dietary_restrictions": {
      "SS": ["not-vegetarian", "not-vegan", "contains-alcohol"]
    },
    "ingredients": {
      "L": [
        {
          "M": {
            "category": {
              "S": "pasta"
            },
            "items": {
              "L": [
                {
                  "M": {
                    "item": {
                      "S": "fresh lasagna sheets"
                    },
                    "amount": {
                      "S": "2 lbs"
                    },
                    "notes": {
                      "S": "preferably homemade"
                    }
                  }
                }
              ]
            }
          }
        },
        {
          "M": {
            "category": {
              "S": "meat_sauce"
            },
            "items": {
              "L": [
                {
                  "M": {
                    "item": {
                      "S": "ground beef"
                    },
                    "amount": {
                      "S": "1.5 lbs"
                    },
                    "quality": {
                      "S": "80/20 blend"
                    }
                  }
                },
                {
                  "M": {
                    "item": {
                      "S": "ground pork"
                    },
                    "amount": {
                      "S": "0.5 lbs"
                    }
                  }
                },
                {
                  "M": {
                    "item": {
                      "S": "pancetta"
                    },
                    "amount": {
                      "S": "4 oz diced"
                    }
                  }
                },
                {
                  "M": {
                    "item": {
                      "S": "san marzano tomatoes"
                    },
                    "amount": {
                      "S": "28 oz can"
                    },
                    "brand": {
                      "S": "imported"
                    }
                  }
                },
                {
                  "M": {
                    "item": {
                      "S": "red wine"
                    },
                    "amount": {
                      "S": "1 cup"
                    },
                    "type": {
                      "S": "chianti classico"
                    }
                  }
                }
              ]
            }
          }
        },
        {
          "M": {
            "category": {
              "S": "bechamel"
            },
            "items": {
              "L": [
                {
                  "M": {
                    "item": {
                      "S": "butter"
                    },
                    "amount": {
                      "S": "6 tbsp"
                    },
                    "type": {
                      "S": "european style"
                    }
                  }
                },
                {
                  "M": {
                    "item": {
                      "S": "all-purpose flour"
                    },
                    "amount": {
                      "S": "6 tbsp"
                    }
                  }
                },
                {
                  "M": {
                    "item": {
                      "S": "whole milk"
                    },
                    "amount": {
                      "S": "4 cups"
                    },
                    "temperature": {
                      "S": "warm"
                    }
                  }
                },
                {
                  "M": {
                    "item": {
                      "S": "nutmeg"
                    },
                    "amount": {
                      "S": "pinch"
                    },
                    "type": {
                      "S": "freshly grated"
                    }
                  }
                }
              ]
            }
          }
        },
        {
          "M": {
            "category": {
              "S": "cheeses"
            },
            "items": {
              "L": [
                {
                  "M": {
                    "item": {
                      "S": "parmigiano-reggiano"
                    },
                    "amount": {
                      "S": "2 cups grated"
                    },
                    "age": {
                      "S": "24 months"
                    }
                  }
                },
                {
                  "M": {
                    "item": {
                      "S": "ricotta"
                    },
                    "amount": {
                      "S": "2 lbs"
                    },
                    "type": {
                      "S": "whole milk"
                    }
                  }
                },
                {
                  "M": {
                    "item": {
                      "S": "mozzarella"
                    },
                    "amount": {
                      "S": "1 lb shredded"
                    },
                    "type": {
                      "S": "low-moisture"
                    }
                  }
                }
              ]
            }
          }
        }
      ]
    },
    "instructions": {
      "L": [
        {
          "M": {
            "step": {
              "N": "1"
            },
            "title": {
              "S": "Prepare Meat Sauce"
            },
            "description": {
              "S": "Brown pancetta, add ground meats, cook with vegetables and wine"
            },
            "time": {
              "N": "45"
            },
            "temperature": {
              "S": "medium-high"
            }
          }
        },
        {
          "M": {
            "step": {
              "N": "2"
            },
            "title": {
              "S": "Make Bechamel"
            },
            "description": {
              "S": "Create roux with butter and flour, gradually add warm milk"
            },
            "time": {
              "N": "20"
            },
            "tips": {
              "L": [
                {
                  "S": "Whisk constantly to prevent lumps"
                },
                {
                  "S": "Keep milk warm for smooth incorporation"
                }
              ]
            }
          }
        },
        {
          "M": {
            "step": {
              "N": "3"
            },
            "title": {
              "S": "Layer Assembly"
            },
            "description": {
              "S": "Alternate layers of pasta, meat sauce, bechamel, and cheeses"
            },
            "time": {
              "N": "30"
            },
            "layers": {
              "L": [
                {
                  "M": {
                    "order": {
                      "N": "1"
                    },
                    "components": {
                      "SS": ["meat_sauce", "pasta", "bechamel", "ricotta"]
                    }
                  }
                },
                {
                  "M": {
                    "order": {
                      "N": "2"
                    },
                    "components": {
                      "SS": ["pasta", "meat_sauce", "bechamel", "mozzarella"]
                    }
                  }
                },
                {
                  "M": {
                    "order": {
                      "N": "3"
                    },
                    "components": {
                      "SS": ["pasta", "meat_sauce", "bechamel", "parmigiano"]
                    }
                  }
                }
              ]
            }
          }
        }
      ]
    },
    "nutrition": {
      "M": {
        "per_serving": {
          "M": {
            "calories": {
              "N": "680"
            },
            "protein": {
              "N": "42"
            },
            "carbohydrates": {
              "N": "35"
            },
            "fat": {
              "N": "38"
            },
            "fiber": {
              "N": "3"
            },
            "sodium": {
              "N": "1250"
            },
            "cholesterol": {
              "N": "145"
            }
          }
        },
        "daily_values": {
          "M": {
            "protein": {
              "N": "84"
            },
            "vitamin_a": {
              "N": "25"
            },
            "calcium": {
              "N": "45"
            },
            "iron": {
              "N": "20"
            }
          }
        }
      }
    },
    "equipment": {
      "L": [
        {
          "M": {
            "item": {
              "S": "9x13 baking dish"
            },
            "essential": {
              "BOOL": true
            }
          }
        },
        {
          "M": {
            "item": {
              "S": "large skillet"
            },
            "essential": {
              "BOOL": true
            }
          }
        },
        {
          "M": {
            "item": {
              "S": "heavy saucepan"
            },
            "essential": {
              "BOOL": true
            }
          }
        },
        {
          "M": {
            "item": {
              "S": "pasta machine"
            },
            "essential": {
              "BOOL": false
            },
            "alternative": {
              "S": "store-bought sheets"
            }
          }
        }
      ]
    },
    "wine_pairing": {
      "M": {
        "primary": {
          "S": "Chianti Classico"
        },
        "alternatives": {
          "SS": ["Sangiovese", "Barbera d'Alba", "Montepulciano"]
        },
        "serving_temp": {
          "S": "60-65°F"
        }
      }
    },
    "storage": {
      "M": {
        "refrigerator": {
          "M": {
            "duration": {
              "S": "3-4 days"
            },
            "container": {
              "S": "covered tightly"
            }
          }
        },
        "freezer": {
          "M": {
            "duration": {
              "S": "3 months"
            },
            "instructions": {
              "L": [
                {
                  "S": "Cool completely before freezing"
                },
                {
                  "S": "Wrap in plastic then foil"
                },
                {
                  "S": "Thaw overnight in refrigerator"
                }
              ]
            }
          }
        }
      }
    },
    "reviews": {
      "L": [
        {
          "M": {
            "rating": {
              "N": "5"
            },
            "comment": {
              "S": "Absolutely incredible! Worth every minute of preparation time."
            },
            "reviewer": {
              "S": "chef_mario_2021"
            },
            "date": {
              "S": "2021-12-15"
            },
            "verified": {
              "BOOL": true
            },
            "helpful_votes": {
              "N": "47"
            }
          }
        },
        {
          "M": {
            "rating": {
              "N": "5"
            },
            "comment": {
              "S": "Family recipe perfection. Made this for Christmas dinner and everyone asked for the recipe!"
            },
            "reviewer": {
              "S": "nonna_rosa"
            },
            "date": {
              "S": "2021-12-25"
            },
            "verified": {
              "BOOL": true
            },
            "helpful_votes": {
              "N": "32"
            }
          }
        },
        {
          "M": {
            "rating": {
              "N": "4"
            },
            "comment": {
              "S": "Delicious but very time consuming. Plan ahead!"
            },
            "reviewer": {
              "S": "busy_parent_123"
            },
            "date": {
              "S": "2021-11-28"
            },
            "verified": {
              "BOOL": true
            },
            "helpful_votes": {
              "N": "18"
            }
          }
        }
      ]
    },
    "recipe_history": {
      "M": {
        "origin": {
          "S": "Emilia-Romagna, Italy"
        },
        "family_generations": {
          "N": "4"
        },
        "first_recorded": {
          "S": "1923"
        },
        "modifications": {
          "L": [
            {
              "M": {
                "year": {
                  "S": "1965"
                },
                "change": {
                  "S": "Added wine to meat sauce"
                },
                "reason": {
                  "S": "Enhanced flavor depth"
                }
              }
            },
            {
              "M": {
                "year": {
                  "S": "1987"
                },
                "change": {
                  "S": "Increased cheese blend variety"
                },
                "reason": {
                  "S": "Improved texture and taste"
                }
              }
            }
          ]
        }
      }
    },
    "cooking_tips": {
      "L": [
        {
          "M": {
            "category": {
              "S": "preparation"
            },
            "tip": {
              "S": "Make sauce day before for better flavor development"
            },
            "importance": {
              "S": "high"
            }
          }
        },
        {
          "M": {
            "category": {
              "S": "assembly"
            },
            "tip": {
              "S": "Let each layer cool slightly before adding the next"
            },
            "importance": {
              "S": "medium"
            }
          }
        },
        {
          "M": {
            "category": {
              "S": "baking"
            },
            "tip": {
              "S": "Cover with foil for first hour, then uncover to brown"
            },
            "importance": {
              "S": "high"
            }
          }
        }
      ]
    }
  },
  "ConsumedCapacity": {
    "TableName": "pasta-recipes",
    "CapacityUnits": 8.5,
    "ReadCapacityUnits": 8.5,
    "WriteCapacityUnits": 0.0
  }
}
"""
        tags: ["serde-benchmark"]
    }
    {
        id: "awsJson1_0_GetItemOutputBinary_S"
        protocol: awsJson1_0
        code: 200
        body: """
{
  "Item": {
    "id": {
      "B": "cmVjaXBlLTAwMQ=="
    },
    "name": {
      "B": "Q2xhc3NpYyBDYXJib25hcmE="
    },
    "cuisine": {
      "B": "SXRhbGlhbg=="
    },
    "cook_time": {
      "N": "20"
    },
    "difficulty": {
      "B": "TWVkaXVt"
    },
    "rating": {
      "N": "4.8"
    }
  },
  "ConsumedCapacity": {
    "TableName": "pasta-recipes",
    "CapacityUnits": 1
  }
}
"""
        tags: ["serde-benchmark"]
    }
    {
        id: "awsJson1_0_GetItemOutputBinary_M"
        protocol: awsJson1_0
        code: 200
        body: """
{
  "Item": {
    "id": {
      "B": "cmVjaXBlLTAwMg=="
    },
    "name": {
      "B": "RmV0dHVjY2luZSBBbGZyZWRv"
    },
    "description": {
      "B": "Q3JlYW15LCByaWNoIHBhc3RhIGRpc2ggd2l0aCBidXR0ZXIsIHBhcm1lc2FuIGNoZWVzZSwgYW5kIGZyZXNoIGZldHR1Y2NpbmUgbm9vZGxlcw=="
    },
    "cook_time": {
      "N": "25"
    },
    "prep_time": {
      "N": "15"
    },
    "difficulty": {
      "B": "RWFzeQ=="
    },
    "cuisine": {
      "B": "SXRhbGlhbg=="
    },
    "servings": {
      "N": "4"
    },
    "rating": {
      "N": "4.6"
    },
    "tags": {
      "BS": [
        "Y3JlYW15",
        "Y29tZm9ydC1mb29k",
        "dmVnZXRhcmlhbg=="
      ]
    },
    "ingredients": {
      "L": [
        {
          "M": {
            "item": {
              "B": "ZmV0dHVjY2luZSBwYXN0YQ=="
            },
            "amount": {
              "B": "MSBsYg=="
            }
          }
        },
        {
          "M": {
            "item": {
              "B": "YnV0dGVy"
            },
            "amount": {
              "B": "MS8yIGN1cA=="
            }
          }
        },
        {
          "M": {
            "item": {
              "B": "cGFybWVzYW4gY2hlZXNl"
            },
            "amount": {
              "B": "MSBjdXAgZ3JhdGVk"
            }
          }
        },
        {
          "M": {
            "item": {
              "B": "aGVhdnkgY3JlYW0="
            },
            "amount": {
              "B": "MS8yIGN1cA=="
            }
          }
        }
      ]
    },
    "nutrition": {
      "M": {
        "calories": {
          "N": "520"
        },
        "protein": {
          "N": "18"
        },
        "carbs": {
          "N": "45"
        },
        "fat": {
          "N": "28"
        }
      }
    }
  },
  "ConsumedCapacity": {
    "TableName": "pasta-recipes",
    "CapacityUnits": 2.5,
    "ReadCapacityUnits": 2.5
  }
}
"""
        tags: ["serde-benchmark"]
    }
    {
        id: "awsJson1_0_GetItemOutputBinary_L"
        protocol: awsJson1_0
        code: 200
        body: """
{
  "Item": {
    "id": {
      "B": "cmVjaXBlLTAwMw=="
    },
    "name": {
      "B": "R3JhbmRtYSdzIFVsdGltYXRlIExhc2FnbmEgQm9sb2duZXNl"
    },
    "description": {
      "B": "QSB0cmFkaXRpb25hbCBJdGFsaWFuIGxhc2FnbmEgcmVjaXBlIHBhc3NlZCBkb3duIHRocm91Z2ggZ2VuZXJhdGlvbnMsIGZlYXR1cmluZyBsYXllcnMgb2YgcmljaCBtZWF0IHNhdWNlLCBjcmVhbXkgYmVjaGFtZWwsIGZyZXNoIHBhc3RhIHNoZWV0cywgYW5kIGEgYmxlbmQgb2YgYXJ0aXNhbmFsIGNoZWVzZXMuIFRoaXMgY29tcGxleCBkaXNoIHJlcXVpcmVzIG11bHRpcGxlIHByZXBhcmF0aW9uIHN0YWdlcyBhbmQgcmVwcmVzZW50cyB0aGUgcGlubmFjbGUgb2YgSXRhbGlhbiBjb21mb3J0IGZvb2QgY3JhZnRzbWFuc2hpcC4gUmVjaXBlIGFkYXB0ZWQgZnJvbSAnTGEgQ3VjaW5hIGRlbGxhIE5vbm5hJyBieSBNYXJpYSBCZW5lZGV0dGksIDE5NTIu"
    },
    "cook_time": {
      "N": "180"
    },
    "prep_time": {
      "N": "120"
    },
    "total_time": {
      "N": "300"
    },
    "difficulty": {
      "B": "RXhwZXJ0"
    },
    "cuisine": {
      "B": "SXRhbGlhbg=="
    },
    "servings": {
      "N": "12"
    },
    "rating": {
      "N": "4.9"
    },
    "cost_estimate": {
      "N": "45.50"
    },
    "active": {
      "BOOL": true
    },
    "featured": {
      "BOOL": true
    },
    "tags": {
      "BS": [
        "dHJhZGl0aW9uYWw=",
        "Y29tZm9ydC1mb29k",
        "ZmFtaWx5LXJlY2lwZQ==",
        "aG9saWRheQ==",
        "bWVhdC1zYXVjZQ==",
        "bGF5ZXJlZA==",
        "YmFrZWQ=",
        "aXRhbGlhbi1jbGFzc2lj",
        "dGltZS1pbnRlbnNpdmU=",
        "c3BlY2lhbC1vY2Nhc2lvbg=="
      ]
    },
    "categories": {
      "BS": [
        "bWFpbi1jb3Vyc2U=",
        "cGFzdGE=",
        "Y2Fzc2Vyb2xl",
        "aXRhbGlhbg=="
      ]
    },
    "allergens": {
      "BS": [
        "ZGFpcnk=",
        "Z2x1dGVu",
        "ZWdncw=="
      ]
    },
    "dietary_restrictions": {
      "BS": [
        "bm90LXZlZ2V0YXJpYW4=",
        "bm90LXZlZ2Fu",
        "Y29udGFpbnMtYWxjb2hvbA=="
      ]
    },
    "ingredients": {
      "L": [
        {
          "M": {
            "category": {
              "B": "cGFzdGE="
            },
            "items": {
              "L": [
                {
                  "M": {
                    "item": {
                      "B": "ZnJlc2ggbGFzYWduYSBzaGVldHM="
                    },
                    "amount": {
                      "B": "MiBsYnM="
                    },
                    "notes": {
                      "B": "cHJlZmVyYWJseSBob21lbWFkZQ=="
                    }
                  }
                }
              ]
            }
          }
        },
        {
          "M": {
            "category": {
              "B": "bWVhdF9zYXVjZQ=="
            },
            "items": {
              "L": [
                {
                  "M": {
                    "item": {
                      "B": "Z3JvdW5kIGJlZWY="
                    },
                    "amount": {
                      "B": "MS41IGxicw=="
                    },
                    "quality": {
                      "B": "ODAvMjAgYmxlbmQ="
                    }
                  }
                },
                {
                  "M": {
                    "item": {
                      "B": "Z3JvdW5kIHBvcms="
                    },
                    "amount": {
                      "B": "MC41IGxicw=="
                    }
                  }
                },
                {
                  "M": {
                    "item": {
                      "B": "cGFuY2V0dGE="
                    },
                    "amount": {
                      "B": "NCBveiBkaWNlZA=="
                    }
                  }
                },
                {
                  "M": {
                    "item": {
                      "B": "c2FuIG1hcnphbm8gdG9tYXRvZXM="
                    },
                    "amount": {
                      "B": "Mjggb3ogY2Fu"
                    },
                    "brand": {
                      "B": "aW1wb3J0ZWQ="
                    }
                  }
                },
                {
                  "M": {
                    "item": {
                      "B": "cmVkIHdpbmU="
                    },
                    "amount": {
                      "B": "MSBjdXA="
                    },
                    "type": {
                      "B": "Y2hpYW50aSBjbGFzc2ljbw=="
                    }
                  }
                }
              ]
            }
          }
        },
        {
          "M": {
            "category": {
              "B": "YmVjaGFtZWw="
            },
            "items": {
              "L": [
                {
                  "M": {
                    "item": {
                      "B": "YnV0dGVy"
                    },
                    "amount": {
                      "B": "NiB0YnNw"
                    },
                    "type": {
                      "B": "ZXVyb3BlYW4gc3R5bGU="
                    }
                  }
                },
                {
                  "M": {
                    "item": {
                      "B": "YWxsLXB1cnBvc2UgZmxvdXI="
                    },
                    "amount": {
                      "B": "NiB0YnNw"
                    }
                  }
                },
                {
                  "M": {
                    "item": {
                      "B": "d2hvbGUgbWlsaw=="
                    },
                    "amount": {
                      "B": "NCBjdXBz"
                    },
                    "temperature": {
                      "B": "d2FybQ=="
                    }
                  }
                },
                {
                  "M": {
                    "item": {
                      "B": "bnV0bWVn"
                    },
                    "amount": {
                      "B": "cGluY2g="
                    },
                    "type": {
                      "B": "ZnJlc2hseSBncmF0ZWQ="
                    }
                  }
                }
              ]
            }
          }
        },
        {
          "M": {
            "category": {
              "B": "Y2hlZXNlcw=="
            },
            "items": {
              "L": [
                {
                  "M": {
                    "item": {
                      "B": "cGFybWlnaWFuby1yZWdnaWFubw=="
                    },
                    "amount": {
                      "B": "MiBjdXBzIGdyYXRlZA=="
                    },
                    "age": {
                      "B": "MjQgbW9udGhz"
                    }
                  }
                },
                {
                  "M": {
                    "item": {
                      "B": "cmljb3R0YQ=="
                    },
                    "amount": {
                      "B": "MiBsYnM="
                    },
                    "type": {
                      "B": "d2hvbGUgbWlsaw=="
                    }
                  }
                },
                {
                  "M": {
                    "item": {
                      "B": "bW96emFyZWxsYQ=="
                    },
                    "amount": {
                      "B": "MSBsYiBzaHJlZGRlZA=="
                    },
                    "type": {
                      "B": "bG93LW1vaXN0dXJl"
                    }
                  }
                }
              ]
            }
          }
        }
      ]
    },
    "instructions": {
      "L": [
        {
          "M": {
            "step": {
              "N": "1"
            },
            "title": {
              "B": "UHJlcGFyZSBNZWF0IFNhdWNl"
            },
            "description": {
              "B": "QnJvd24gcGFuY2V0dGEsIGFkZCBncm91bmQgbWVhdHMsIGNvb2sgd2l0aCB2ZWdldGFibGVzIGFuZCB3aW5l"
            },
            "time": {
              "N": "45"
            },
            "temperature": {
              "B": "bWVkaXVtLWhpZ2g="
            }
          }
        },
        {
          "M": {
            "step": {
              "N": "2"
            },
            "title": {
              "B": "TWFrZSBCZWNoYW1lbA=="
            },
            "description": {
              "B": "Q3JlYXRlIHJvdXggd2l0aCBidXR0ZXIgYW5kIGZsb3VyLCBncmFkdWFsbHkgYWRkIHdhcm0gbWlsaw=="
            },
            "time": {
              "N": "20"
            },
            "tips": {
              "L": [
                {
                  "B": "V2hpc2sgY29uc3RhbnRseSB0byBwcmV2ZW50IGx1bXBz"
                },
                {
                  "B": "S2VlcCBtaWxrIHdhcm0gZm9yIHNtb290aCBpbmNvcnBvcmF0aW9u"
                }
              ]
            }
          }
        },
        {
          "M": {
            "step": {
              "N": "3"
            },
            "title": {
              "B": "TGF5ZXIgQXNzZW1ibHk="
            },
            "description": {
              "B": "QWx0ZXJuYXRlIGxheWVycyBvZiBwYXN0YSwgbWVhdCBzYXVjZSwgYmVjaGFtZWwsIGFuZCBjaGVlc2Vz"
            },
            "time": {
              "N": "30"
            },
            "layers": {
              "L": [
                {
                  "M": {
                    "order": {
                      "N": "1"
                    },
                    "components": {
                      "BS": [
                        "bWVhdF9zYXVjZQ==",
                        "cGFzdGE=",
                        "YmVjaGFtZWw=",
                        "cmljb3R0YQ=="
                      ]
                    }
                  }
                },
                {
                  "M": {
                    "order": {
                      "N": "2"
                    },
                    "components": {
                      "BS": [
                        "cGFzdGE=",
                        "bWVhdF9zYXVjZQ==",
                        "YmVjaGFtZWw=",
                        "bW96emFyZWxsYQ=="
                      ]
                    }
                  }
                },
                {
                  "M": {
                    "order": {
                      "N": "3"
                    },
                    "components": {
                      "BS": [
                        "cGFzdGE=",
                        "bWVhdF9zYXVjZQ==",
                        "YmVjaGFtZWw=",
                        "cGFybWlnaWFubw=="
                      ]
                    }
                  }
                }
              ]
            }
          }
        }
      ]
    },
    "nutrition": {
      "M": {
        "per_serving": {
          "M": {
            "calories": {
              "N": "680"
            },
            "protein": {
              "N": "42"
            },
            "carbohydrates": {
              "N": "35"
            },
            "fat": {
              "N": "38"
            },
            "fiber": {
              "N": "3"
            },
            "sodium": {
              "N": "1250"
            },
            "cholesterol": {
              "N": "145"
            }
          }
        },
        "daily_values": {
          "M": {
            "protein": {
              "N": "84"
            },
            "vitamin_a": {
              "N": "25"
            },
            "calcium": {
              "N": "45"
            },
            "iron": {
              "N": "20"
            }
          }
        }
      }
    },
    "equipment": {
      "L": [
        {
          "M": {
            "item": {
              "B": "OXgxMyBiYWtpbmcgZGlzaA=="
            },
            "essential": {
              "BOOL": true
            }
          }
        },
        {
          "M": {
            "item": {
              "B": "bGFyZ2Ugc2tpbGxldA=="
            },
            "essential": {
              "BOOL": true
            }
          }
        },
        {
          "M": {
            "item": {
              "B": "aGVhdnkgc2F1Y2VwYW4="
            },
            "essential": {
              "BOOL": true
            }
          }
        },
        {
          "M": {
            "item": {
              "B": "cGFzdGEgbWFjaGluZQ=="
            },
            "essential": {
              "BOOL": false
            },
            "alternative": {
              "B": "c3RvcmUtYm91Z2h0IHNoZWV0cw=="
            }
          }
        }
      ]
    },
    "wine_pairing": {
      "M": {
        "primary": {
          "B": "Q2hpYW50aSBDbGFzc2ljbw=="
        },
        "alternatives": {
          "BS": [
            "U2FuZ2lvdmVzZQ==",
            "QmFyYmVyYSBkJ0FsYmE=",
            "TW9udGVwdWxjaWFubw=="
          ]
        },
        "serving_temp": {
          "B": "NjAtNjXCsEY="
        }
      }
    },
    "storage": {
      "M": {
        "refrigerator": {
          "M": {
            "duration": {
              "B": "My00IGRheXM="
            },
            "container": {
              "B": "Y292ZXJlZCB0aWdodGx5"
            }
          }
        },
        "freezer": {
          "M": {
            "duration": {
              "B": "MyBtb250aHM="
            },
            "instructions": {
              "L": [
                {
                  "B": "Q29vbCBjb21wbGV0ZWx5IGJlZm9yZSBmcmVlemluZw=="
                },
                {
                  "B": "V3JhcCBpbiBwbGFzdGljIHRoZW4gZm9pbA=="
                },
                {
                  "B": "VGhhdyBvdmVybmlnaHQgaW4gcmVmcmlnZXJhdG9y"
                }
              ]
            }
          }
        }
      }
    },
    "reviews": {
      "L": [
        {
          "M": {
            "rating": {
              "N": "5"
            },
            "comment": {
              "B": "QWJzb2x1dGVseSBpbmNyZWRpYmxlISBXb3J0aCBldmVyeSBtaW51dGUgb2YgcHJlcGFyYXRpb24gdGltZS4="
            },
            "reviewer": {
              "B": "Y2hlZl9tYXJpb18yMDIx"
            },
            "date": {
              "B": "MjAyMS0xMi0xNQ=="
            },
            "verified": {
              "BOOL": true
            },
            "helpful_votes": {
              "N": "47"
            }
          }
        },
        {
          "M": {
            "rating": {
              "N": "5"
            },
            "comment": {
              "B": "RmFtaWx5IHJlY2lwZSBwZXJmZWN0aW9uLiBNYWRlIHRoaXMgZm9yIENocmlzdG1hcyBkaW5uZXIgYW5kIGV2ZXJ5b25lIGFza2VkIGZvciB0aGUgcmVjaXBlIQ=="
            },
            "reviewer": {
              "B": "bm9ubmFfcm9zYQ=="
            },
            "date": {
              "B": "MjAyMS0xMi0yNQ=="
            },
            "verified": {
              "BOOL": true
            },
            "helpful_votes": {
              "N": "32"
            }
          }
        },
        {
          "M": {
            "rating": {
              "N": "4"
            },
            "comment": {
              "B": "RGVsaWNpb3VzIGJ1dCB2ZXJ5IHRpbWUgY29uc3VtaW5nLiBQbGFuIGFoZWFkIQ=="
            },
            "reviewer": {
              "B": "YnVzeV9wYXJlbnRfMTIz"
            },
            "date": {
              "B": "MjAyMS0xMS0yOA=="
            },
            "verified": {
              "BOOL": true
            },
            "helpful_votes": {
              "N": "18"
            }
          }
        }
      ]
    },
    "recipe_history": {
      "M": {
        "origin": {
          "B": "RW1pbGlhLVJvbWFnbmEsIEl0YWx5"
        },
        "family_generations": {
          "N": "4"
        },
        "first_recorded": {
          "B": "MTkyMw=="
        },
        "modifications": {
          "L": [
            {
              "M": {
                "year": {
                  "B": "MTk2NQ=="
                },
                "change": {
                  "B": "QWRkZWQgd2luZSB0byBtZWF0IHNhdWNl"
                },
                "reason": {
                  "B": "RW5oYW5jZWQgZmxhdm9yIGRlcHRo"
                }
              }
            },
            {
              "M": {
                "year": {
                  "B": "MTk4Nw=="
                },
                "change": {
                  "B": "SW5jcmVhc2VkIGNoZWVzZSBibGVuZCB2YXJpZXR5"
                },
                "reason": {
                  "B": "SW1wcm92ZWQgdGV4dHVyZSBhbmQgdGFzdGU="
                }
              }
            }
          ]
        }
      }
    },
    "cooking_tips": {
      "L": [
        {
          "M": {
            "category": {
              "B": "cHJlcGFyYXRpb24="
            },
            "tip": {
              "B": "TWFrZSBzYXVjZSBkYXkgYmVmb3JlIGZvciBiZXR0ZXIgZmxhdm9yIGRldmVsb3BtZW50"
            },
            "importance": {
              "B": "aGlnaA=="
            }
          }
        },
        {
          "M": {
            "category": {
              "B": "YXNzZW1ibHk="
            },
            "tip": {
              "B": "TGV0IGVhY2ggbGF5ZXIgY29vbCBzbGlnaHRseSBiZWZvcmUgYWRkaW5nIHRoZSBuZXh0"
            },
            "importance": {
              "B": "bWVkaXVt"
            }
          }
        },
        {
          "M": {
            "category": {
              "B": "YmFraW5n"
            },
            "tip": {
              "B": "Q292ZXIgd2l0aCBmb2lsIGZvciBmaXJzdCBob3VyLCB0aGVuIHVuY292ZXIgdG8gYnJvd24="
            },
            "importance": {
              "B": "aGlnaA=="
            }
          }
        }
      ]
    }
  },
  "ConsumedCapacity": {
    "TableName": "pasta-recipes",
    "CapacityUnits": 8.5,
    "ReadCapacityUnits": 8.5,
    "WriteCapacityUnits": 0
  }
}
"""
        tags: ["serde-benchmark"]
    }
    // section: cbor
    {
        id: "rpcv2Cbor_GetItemOutput_Baseline"
        protocol: rpcv2Cbor
        headers: {
            "smithy-protocol": "rpc-v2-cbor"
        }
        body: """
oA==
"""
        code: 200
        tags: ["serde-benchmark"]
    }
    {
        id: "rpcv2Cbor_GetItemOutput_S"
        protocol: rpcv2Cbor
        headers: {
            "smithy-protocol": "rpc-v2-cbor"
        }
        code: 200
        body: """
omRJdGVtpmJpZKFhU2pyZWNpcGUtMDAxZG5hbWWhYVNxQ2xhc3NpYyBDYXJib25hcmFnY3Vpc2luZaFhU2dJdGFsaWFuaWNvb2tfdGltZaFhTmIyMGpkaWZmaWN1bHR5oWFTZk1lZGl1bWZyYXRpbmehYU5jNC44cENvbnN1bWVkQ2FwYWNpdHmiaVRhYmxlTmFtZW1wYXN0YS1yZWNpcGVzbUNhcGFjaXR5VW5pdHP7P/GZmZmZmZo=
"""
        tags: ["serde-benchmark"]
    }
    {
        id: "rpcv2Cbor_GetItemOutput_M"
        protocol: rpcv2Cbor
        headers: {
            "smithy-protocol": "rpc-v2-cbor"
        }
        code: 200
        body: """
omRJdGVtrGJpZKFhU2pyZWNpcGUtMDAyZG5hbWWhYVNyRmV0dHVjY2luZSBBbGZyZWRva2Rlc2NyaXB0aW9uoWFTeFJDcmVhbXksIHJpY2ggcGFzdGEgZGlzaCB3aXRoIGJ1dHRlciwgcGFybWVzYW4gY2hlZXNlLCBhbmQgZnJlc2ggZmV0dHVjY2luZSBub29kbGVzaWNvb2tfdGltZaFhTmIyNWlwcmVwX3RpbWWhYU5iMTVqZGlmZmljdWx0eaFhU2RFYXN5Z2N1aXNpbmWhYVNnSXRhbGlhbmhzZXJ2aW5nc6FhTmE0ZnJhdGluZ6FhTmM0LjZkdGFnc6FiU1ODZmNyZWFteWxjb21mb3J0LWZvb2RqdmVnZXRhcmlhbmtpbmdyZWRpZW50c6FhTIShYU2iZGl0ZW2hYVNwZmV0dHVjY2luZSBwYXN0YWZhbW91bnShYVNkMSBsYqFhTaJkaXRlbaFhU2ZidXR0ZXJmYW1vdW50oWFTZzEvMiBjdXChYU2iZGl0ZW2hYVNvcGFybWVzYW4gY2hlZXNlZmFtb3VudKFhU2wxIGN1cCBncmF0ZWShYU2iZGl0ZW2hYVNraGVhdnkgY3JlYW1mYW1vdW50oWFTZzEvMiBjdXBpbnV0cml0aW9uoWFNpGhjYWxvcmllc6FhTmM1MjBncHJvdGVpbqFhTmIxOGVjYXJic6FhTmI0NWNmYXShYU5iMjhwQ29uc3VtZWRDYXBhY2l0eaNpVGFibGVOYW1lbXBhc3RhLXJlY2lwZXNtQ2FwYWNpdHlVbml0c/tABAAAAAAAAHFSZWFkQ2FwYWNpdHlVbml0c/tABAAAAAAAAA==
"""
        tags: ["serde-benchmark"]
    }
    {
        id: "rpcv2Cbor_GetItemOutput_L"
        protocol: rpcv2Cbor
        headers: {
            "smithy-protocol": "rpc-v2-cbor"
        }
        code: 200
        body: """
omRJdGVtuBpiaWShYVNqcmVjaXBlLTAwM2RuYW1loWFTeCRHcmFuZG1hJ3MgVWx0aW1hdGUgTGFzYWduYSBCb2xvZ25lc2VrZGVzY3JpcHRpb26hYVN5AXFBIHRyYWRpdGlvbmFsIEl0YWxpYW4gbGFzYWduYSByZWNpcGUgcGFzc2VkIGRvd24gdGhyb3VnaCBnZW5lcmF0aW9ucywgZmVhdHVyaW5nIGxheWVycyBvZiByaWNoIG1lYXQgc2F1Y2UsIGNyZWFteSBiZWNoYW1lbCwgZnJlc2ggcGFzdGEgc2hlZXRzLCBhbmQgYSBibGVuZCBvZiBhcnRpc2FuYWwgY2hlZXNlcy4gVGhpcyBjb21wbGV4IGRpc2ggcmVxdWlyZXMgbXVsdGlwbGUgcHJlcGFyYXRpb24gc3RhZ2VzIGFuZCByZXByZXNlbnRzIHRoZSBwaW5uYWNsZSBvZiBJdGFsaWFuIGNvbWZvcnQgZm9vZCBjcmFmdHNtYW5zaGlwLiBSZWNpcGUgYWRhcHRlZCBmcm9tICdMYSBDdWNpbmEgZGVsbGEgTm9ubmEnIGJ5IE1hcmlhIEJlbmVkZXR0aSwgMTk1Mi5pY29va190aW1loWFOYzE4MGlwcmVwX3RpbWWhYU5jMTIwanRvdGFsX3RpbWWhYU5jMzAwamRpZmZpY3VsdHmhYVNmRXhwZXJ0Z2N1aXNpbmWhYVNnSXRhbGlhbmhzZXJ2aW5nc6FhTmIxMmZyYXRpbmehYU5jNC45bWNvc3RfZXN0aW1hdGWhYU5lNDUuNTBmYWN0aXZloWRCT09M9WhmZWF0dXJlZKFkQk9PTPVkdGFnc6FiU1OKa3RyYWRpdGlvbmFsbGNvbWZvcnQtZm9vZG1mYW1pbHktcmVjaXBlZ2hvbGlkYXlqbWVhdC1zYXVjZWdsYXllcmVkZWJha2Vkb2l0YWxpYW4tY2xhc3NpY250aW1lLWludGVuc2l2ZXBzcGVjaWFsLW9jY2FzaW9uamNhdGVnb3JpZXOhYlNThGttYWluLWNvdXJzZWVwYXN0YWljYXNzZXJvbGVnaXRhbGlhbmlhbGxlcmdlbnOhYlNTg2VkYWlyeWZnbHV0ZW5kZWdnc3RkaWV0YXJ5X3Jlc3RyaWN0aW9uc6FiU1ODbm5vdC12ZWdldGFyaWFuaW5vdC12ZWdhbnBjb250YWlucy1hbGNvaG9sa2luZ3JlZGllbnRzoWFMhKFhTaJoY2F0ZWdvcnmhYVNlcGFzdGFlaXRlbXOhYUyBoWFNo2RpdGVtoWFTdGZyZXNoIGxhc2FnbmEgc2hlZXRzZmFtb3VudKFhU2UyIGxic2Vub3Rlc6FhU3NwcmVmZXJhYmx5IGhvbWVtYWRloWFNomhjYXRlZ29yeaFhU2ptZWF0X3NhdWNlZWl0ZW1zoWFMhaFhTaNkaXRlbaFhU2tncm91bmQgYmVlZmZhbW91bnShYVNnMS41IGxic2dxdWFsaXR5oWFTazgwLzIwIGJsZW5koWFNomRpdGVtoWFTa2dyb3VuZCBwb3JrZmFtb3VudKFhU2cwLjUgbGJzoWFNomRpdGVtoWFTaHBhbmNldHRhZmFtb3VudKFhU2o0IG96IGRpY2VkoWFNo2RpdGVtoWFTdHNhbiBtYXJ6YW5vIHRvbWF0b2VzZmFtb3VudKFhU2kyOCBveiBjYW5lYnJhbmShYVNoaW1wb3J0ZWShYU2jZGl0ZW2hYVNocmVkIHdpbmVmYW1vdW50oWFTZTEgY3VwZHR5cGWhYVNwY2hpYW50aSBjbGFzc2ljb6FhTaJoY2F0ZWdvcnmhYVNoYmVjaGFtZWxlaXRlbXOhYUyEoWFNo2RpdGVtoWFTZmJ1dHRlcmZhbW91bnShYVNmNiB0YnNwZHR5cGWhYVNuZXVyb3BlYW4gc3R5bGWhYU2iZGl0ZW2hYVNxYWxsLXB1cnBvc2UgZmxvdXJmYW1vdW50oWFTZjYgdGJzcKFhTaNkaXRlbaFhU2p3aG9sZSBtaWxrZmFtb3VudKFhU2Y0IGN1cHNrdGVtcGVyYXR1cmWhYVNkd2FybaFhTaNkaXRlbaFhU2ZudXRtZWdmYW1vdW50oWFTZXBpbmNoZHR5cGWhYVNuZnJlc2hseSBncmF0ZWShYU2iaGNhdGVnb3J5oWFTZ2NoZWVzZXNlaXRlbXOhYUyDoWFNo2RpdGVtoWFTc3Bhcm1pZ2lhbm8tcmVnZ2lhbm9mYW1vdW50oWFTbTIgY3VwcyBncmF0ZWRjYWdloWFTaTI0IG1vbnRoc6FhTaNkaXRlbaFhU2dyaWNvdHRhZmFtb3VudKFhU2UyIGxic2R0eXBloWFTandob2xlIG1pbGuhYU2jZGl0ZW2hYVNqbW96emFyZWxsYWZhbW91bnShYVNtMSBsYiBzaHJlZGRlZGR0eXBloWFTbGxvdy1tb2lzdHVyZWxpbnN0cnVjdGlvbnOhYUyDoWFNpWRzdGVwoWFOYTFldGl0bGWhYVNyUHJlcGFyZSBNZWF0IFNhdWNla2Rlc2NyaXB0aW9uoWFTeD9Ccm93biBwYW5jZXR0YSwgYWRkIGdyb3VuZCBtZWF0cywgY29vayB3aXRoIHZlZ2V0YWJsZXMgYW5kIHdpbmVkdGltZaFhTmI0NWt0ZW1wZXJhdHVyZaFhU2ttZWRpdW0taGlnaKFhTaVkc3RlcKFhTmEyZXRpdGxloWFTbU1ha2UgQmVjaGFtZWxrZGVzY3JpcHRpb26hYVN4OkNyZWF0ZSByb3V4IHdpdGggYnV0dGVyIGFuZCBmbG91ciwgZ3JhZHVhbGx5IGFkZCB3YXJtIG1pbGtkdGltZaFhTmIyMGR0aXBzoWFMgqFhU3ghV2hpc2sgY29uc3RhbnRseSB0byBwcmV2ZW50IGx1bXBzoWFTeCdLZWVwIG1pbGsgd2FybSBmb3Igc21vb3RoIGluY29ycG9yYXRpb26hYU2lZHN0ZXChYU5hM2V0aXRsZaFhU25MYXllciBBc3NlbWJseWtkZXNjcmlwdGlvbqFhU3g8QWx0ZXJuYXRlIGxheWVycyBvZiBwYXN0YSwgbWVhdCBzYXVjZSwgYmVjaGFtZWwsIGFuZCBjaGVlc2VzZHRpbWWhYU5iMzBmbGF5ZXJzoWFMg6FhTaJlb3JkZXKhYU5hMWpjb21wb25lbnRzoWJTU4RqbWVhdF9zYXVjZWVwYXN0YWhiZWNoYW1lbGdyaWNvdHRhoWFNomVvcmRlcqFhTmEyamNvbXBvbmVudHOhYlNThGVwYXN0YWptZWF0X3NhdWNlaGJlY2hhbWVsam1venphcmVsbGGhYU2iZW9yZGVyoWFOYTNqY29tcG9uZW50c6FiU1OEZXBhc3Rham1lYXRfc2F1Y2VoYmVjaGFtZWxqcGFybWlnaWFub2ludXRyaXRpb26hYU2ia3Blcl9zZXJ2aW5noWFNp2hjYWxvcmllc6FhTmM2ODBncHJvdGVpbqFhTmI0Mm1jYXJib2h5ZHJhdGVzoWFOYjM1Y2ZhdKFhTmIzOGVmaWJlcqFhTmEzZnNvZGl1baFhTmQxMjUwa2Nob2xlc3Rlcm9soWFOYzE0NWxkYWlseV92YWx1ZXOhYU2kZ3Byb3RlaW6hYU5iODRpdml0YW1pbl9hoWFOYjI1Z2NhbGNpdW2hYU5iNDVkaXJvbqFhTmIyMGllcXVpcG1lbnShYUyEoWFNomRpdGVtoWFTcDl4MTMgYmFraW5nIGRpc2hpZXNzZW50aWFsoWRCT09M9aFhTaJkaXRlbaFhU21sYXJnZSBza2lsbGV0aWVzc2VudGlhbKFkQk9PTPWhYU2iZGl0ZW2hYVNuaGVhdnkgc2F1Y2VwYW5pZXNzZW50aWFsoWRCT09M9aFhTaNkaXRlbaFhU21wYXN0YSBtYWNoaW5laWVzc2VudGlhbKFkQk9PTPRrYWx0ZXJuYXRpdmWhYVNzc3RvcmUtYm91Z2h0IHNoZWV0c2x3aW5lX3BhaXJpbmehYU2jZ3ByaW1hcnmhYVNwQ2hpYW50aSBDbGFzc2ljb2xhbHRlcm5hdGl2ZXOhYlNTg2pTYW5naW92ZXNlbkJhcmJlcmEgZCdBbGJhbU1vbnRlcHVsY2lhbm9sc2VydmluZ190ZW1woWFTaDYwLTY1wrBGZ3N0b3JhZ2WhYU2ibHJlZnJpZ2VyYXRvcqFhTaJoZHVyYXRpb26hYVNoMy00IGRheXNpY29udGFpbmVyoWFTb2NvdmVyZWQgdGlnaHRseWdmcmVlemVyoWFNomhkdXJhdGlvbqFhU2gzIG1vbnRoc2xpbnN0cnVjdGlvbnOhYUyDoWFTeB9Db29sIGNvbXBsZXRlbHkgYmVmb3JlIGZyZWV6aW5noWFTeBlXcmFwIGluIHBsYXN0aWMgdGhlbiBmb2lsoWFTeB5UaGF3IG92ZXJuaWdodCBpbiByZWZyaWdlcmF0b3JncmV2aWV3c6FhTIOhYU2mZnJhdGluZ6FhTmE1Z2NvbW1lbnShYVN4PkFic29sdXRlbHkgaW5jcmVkaWJsZSEgV29ydGggZXZlcnkgbWludXRlIG9mIHByZXBhcmF0aW9uIHRpbWUuaHJldmlld2VyoWFTb2NoZWZfbWFyaW9fMjAyMWRkYXRloWFTajIwMjEtMTItMTVodmVyaWZpZWShZEJPT0z1bWhlbHBmdWxfdm90ZXOhYU5iNDehYU2mZnJhdGluZ6FhTmE1Z2NvbW1lbnShYVN4W0ZhbWlseSByZWNpcGUgcGVyZmVjdGlvbi4gTWFkZSB0aGlzIGZvciBDaHJpc3RtYXMgZGlubmVyIGFuZCBldmVyeW9uZSBhc2tlZCBmb3IgdGhlIHJlY2lwZSFocmV2aWV3ZXKhYVNqbm9ubmFfcm9zYWRkYXRloWFTajIwMjEtMTItMjVodmVyaWZpZWShZEJPT0z1bWhlbHBmdWxfdm90ZXOhYU5iMzKhYU2mZnJhdGluZ6FhTmE0Z2NvbW1lbnShYVN4LkRlbGljaW91cyBidXQgdmVyeSB0aW1lIGNvbnN1bWluZy4gUGxhbiBhaGVhZCFocmV2aWV3ZXKhYVNvYnVzeV9wYXJlbnRfMTIzZGRhdGWhYVNqMjAyMS0xMS0yOGh2ZXJpZmllZKFkQk9PTPVtaGVscGZ1bF92b3Rlc6FhTmIxOG5yZWNpcGVfaGlzdG9yeaFhTaRmb3JpZ2luoWFTdUVtaWxpYS1Sb21hZ25hLCBJdGFseXJmYW1pbHlfZ2VuZXJhdGlvbnOhYU5hNG5maXJzdF9yZWNvcmRlZKFhU2QxOTIzbW1vZGlmaWNhdGlvbnOhYUyCoWFNo2R5ZWFyoWFTZDE5NjVmY2hhbmdloWFTeBhBZGRlZCB3aW5lIHRvIG1lYXQgc2F1Y2VmcmVhc29uoWFTdUVuaGFuY2VkIGZsYXZvciBkZXB0aKFhTaNkeWVhcqFhU2QxOTg3ZmNoYW5nZaFhU3geSW5jcmVhc2VkIGNoZWVzZSBibGVuZCB2YXJpZXR5ZnJlYXNvbqFhU3gaSW1wcm92ZWQgdGV4dHVyZSBhbmQgdGFzdGVsY29va2luZ190aXBzoWFMg6FhTaNoY2F0ZWdvcnmhYVNrcHJlcGFyYXRpb25jdGlwoWFTeDNNYWtlIHNhdWNlIGRheSBiZWZvcmUgZm9yIGJldHRlciBmbGF2b3IgZGV2ZWxvcG1lbnRqaW1wb3J0YW5jZaFhU2RoaWdooWFNo2hjYXRlZ29yeaFhU2hhc3NlbWJseWN0aXChYVN4M0xldCBlYWNoIGxheWVyIGNvb2wgc2xpZ2h0bHkgYmVmb3JlIGFkZGluZyB0aGUgbmV4dGppbXBvcnRhbmNloWFTZm1lZGl1baFhTaNoY2F0ZWdvcnmhYVNmYmFraW5nY3RpcKFhU3g1Q292ZXIgd2l0aCBmb2lsIGZvciBmaXJzdCBob3VyLCB0aGVuIHVuY292ZXIgdG8gYnJvd25qaW1wb3J0YW5jZaFhU2RoaWdocENvbnN1bWVkQ2FwYWNpdHmkaVRhYmxlTmFtZW1wYXN0YS1yZWNpcGVzbUNhcGFjaXR5VW5pdHP7QCEAAAAAAABxUmVhZENhcGFjaXR5VW5pdHP7QCEAAAAAAAByV3JpdGVDYXBhY2l0eVVuaXRzAA==
"""
        tags: ["serde-benchmark"]
    }
    {
        id: "rpcv2Cbor_GetItemOutputBinary_S"
        protocol: rpcv2Cbor
        headers: {
            "smithy-protocol": "rpc-v2-cbor"
        }
        code: 200
        body: """
omRJdGVtpmJpZKFhQkpyZWNpcGUtMDAxZG5hbWWhYUJRQ2xhc3NpYyBDYXJib25hcmFnY3Vpc2luZaFhQkdJdGFsaWFuaWNvb2tfdGltZaFhTmIyMGpkaWZmaWN1bHR5oWFCRk1lZGl1bWZyYXRpbmehYU5jNC44cENvbnN1bWVkQ2FwYWNpdHmiaVRhYmxlTmFtZW1wYXN0YS1yZWNpcGVzbUNhcGFjaXR5VW5pdHMB
"""
        tags: ["serde-benchmark"]
    }
    {
        id: "rpcv2Cbor_GetItemOutputBinary_M"
        protocol: rpcv2Cbor
        headers: {
            "smithy-protocol": "rpc-v2-cbor"
        }
        code: 200
        body: """
omRJdGVtrGJpZKFhQkpyZWNpcGUtMDAyZG5hbWWhYUJSRmV0dHVjY2luZSBBbGZyZWRva2Rlc2NyaXB0aW9uoWFCWFJDcmVhbXksIHJpY2ggcGFzdGEgZGlzaCB3aXRoIGJ1dHRlciwgcGFybWVzYW4gY2hlZXNlLCBhbmQgZnJlc2ggZmV0dHVjY2luZSBub29kbGVzaWNvb2tfdGltZaFhTmIyNWlwcmVwX3RpbWWhYU5iMTVqZGlmZmljdWx0eaFhQkRFYXN5Z2N1aXNpbmWhYUJHSXRhbGlhbmhzZXJ2aW5nc6FhTmE0ZnJhdGluZ6FhTmM0LjZkdGFnc6FiQlODRmNyZWFteUxjb21mb3J0LWZvb2RKdmVnZXRhcmlhbmtpbmdyZWRpZW50c6FhTIShYU2iZGl0ZW2hYUJQZmV0dHVjY2luZSBwYXN0YWZhbW91bnShYUJEMSBsYqFhTaJkaXRlbaFhQkZidXR0ZXJmYW1vdW50oWFCRzEvMiBjdXChYU2iZGl0ZW2hYUJPcGFybWVzYW4gY2hlZXNlZmFtb3VudKFhQkwxIGN1cCBncmF0ZWShYU2iZGl0ZW2hYUJLaGVhdnkgY3JlYW1mYW1vdW50oWFCRzEvMiBjdXBpbnV0cml0aW9uoWFNpGhjYWxvcmllc6FhTmM1MjBncHJvdGVpbqFhTmIxOGVjYXJic6FhTmI0NWNmYXShYU5iMjhwQ29uc3VtZWRDYXBhY2l0eaNpVGFibGVOYW1lbXBhc3RhLXJlY2lwZXNtQ2FwYWNpdHlVbml0c/tABAAAAAAAAHFSZWFkQ2FwYWNpdHlVbml0c/tABAAAAAAAAA==
"""
        tags: ["serde-benchmark"]
    }
    {
        id: "rpcv2Cbor_GetItemOutputBinary_L"
        protocol: rpcv2Cbor
        headers: {
            "smithy-protocol": "rpc-v2-cbor"
        }
        code: 200
        body: """
omRJdGVtuBpiaWShYUJKcmVjaXBlLTAwM2RuYW1loWFCWCRHcmFuZG1hJ3MgVWx0aW1hdGUgTGFzYWduYSBCb2xvZ25lc2VrZGVzY3JpcHRpb26hYUJZAXFBIHRyYWRpdGlvbmFsIEl0YWxpYW4gbGFzYWduYSByZWNpcGUgcGFzc2VkIGRvd24gdGhyb3VnaCBnZW5lcmF0aW9ucywgZmVhdHVyaW5nIGxheWVycyBvZiByaWNoIG1lYXQgc2F1Y2UsIGNyZWFteSBiZWNoYW1lbCwgZnJlc2ggcGFzdGEgc2hlZXRzLCBhbmQgYSBibGVuZCBvZiBhcnRpc2FuYWwgY2hlZXNlcy4gVGhpcyBjb21wbGV4IGRpc2ggcmVxdWlyZXMgbXVsdGlwbGUgcHJlcGFyYXRpb24gc3RhZ2VzIGFuZCByZXByZXNlbnRzIHRoZSBwaW5uYWNsZSBvZiBJdGFsaWFuIGNvbWZvcnQgZm9vZCBjcmFmdHNtYW5zaGlwLiBSZWNpcGUgYWRhcHRlZCBmcm9tICdMYSBDdWNpbmEgZGVsbGEgTm9ubmEnIGJ5IE1hcmlhIEJlbmVkZXR0aSwgMTk1Mi5pY29va190aW1loWFOYzE4MGlwcmVwX3RpbWWhYU5jMTIwanRvdGFsX3RpbWWhYU5jMzAwamRpZmZpY3VsdHmhYUJGRXhwZXJ0Z2N1aXNpbmWhYUJHSXRhbGlhbmhzZXJ2aW5nc6FhTmIxMmZyYXRpbmehYU5jNC45bWNvc3RfZXN0aW1hdGWhYU5lNDUuNTBmYWN0aXZloWRCT09M9WhmZWF0dXJlZKFkQk9PTPVkdGFnc6FiQlOKS3RyYWRpdGlvbmFsTGNvbWZvcnQtZm9vZE1mYW1pbHktcmVjaXBlR2hvbGlkYXlKbWVhdC1zYXVjZUdsYXllcmVkRWJha2VkT2l0YWxpYW4tY2xhc3NpY050aW1lLWludGVuc2l2ZVBzcGVjaWFsLW9jY2FzaW9uamNhdGVnb3JpZXOhYkJThEttYWluLWNvdXJzZUVwYXN0YUljYXNzZXJvbGVHaXRhbGlhbmlhbGxlcmdlbnOhYkJTg0VkYWlyeUZnbHV0ZW5EZWdnc3RkaWV0YXJ5X3Jlc3RyaWN0aW9uc6FiQlODTm5vdC12ZWdldGFyaWFuSW5vdC12ZWdhblBjb250YWlucy1hbGNvaG9sa2luZ3JlZGllbnRzoWFMhKFhTaJoY2F0ZWdvcnmhYUJFcGFzdGFlaXRlbXOhYUyBoWFNo2RpdGVtoWFCVGZyZXNoIGxhc2FnbmEgc2hlZXRzZmFtb3VudKFhQkUyIGxic2Vub3Rlc6FhQlNwcmVmZXJhYmx5IGhvbWVtYWRloWFNomhjYXRlZ29yeaFhQkptZWF0X3NhdWNlZWl0ZW1zoWFMhaFhTaNkaXRlbaFhQktncm91bmQgYmVlZmZhbW91bnShYUJHMS41IGxic2dxdWFsaXR5oWFCSzgwLzIwIGJsZW5koWFNomRpdGVtoWFCS2dyb3VuZCBwb3JrZmFtb3VudKFhQkcwLjUgbGJzoWFNomRpdGVtoWFCSHBhbmNldHRhZmFtb3VudKFhQko0IG96IGRpY2VkoWFNo2RpdGVtoWFCVHNhbiBtYXJ6YW5vIHRvbWF0b2VzZmFtb3VudKFhQkkyOCBveiBjYW5lYnJhbmShYUJIaW1wb3J0ZWShYU2jZGl0ZW2hYUJIcmVkIHdpbmVmYW1vdW50oWFCRTEgY3VwZHR5cGWhYUJQY2hpYW50aSBjbGFzc2ljb6FhTaJoY2F0ZWdvcnmhYUJIYmVjaGFtZWxlaXRlbXOhYUyEoWFNo2RpdGVtoWFCRmJ1dHRlcmZhbW91bnShYUJGNiB0YnNwZHR5cGWhYUJOZXVyb3BlYW4gc3R5bGWhYU2iZGl0ZW2hYUJRYWxsLXB1cnBvc2UgZmxvdXJmYW1vdW50oWFCRjYgdGJzcKFhTaNkaXRlbaFhQkp3aG9sZSBtaWxrZmFtb3VudKFhQkY0IGN1cHNrdGVtcGVyYXR1cmWhYUJEd2FybaFhTaNkaXRlbaFhQkZudXRtZWdmYW1vdW50oWFCRXBpbmNoZHR5cGWhYUJOZnJlc2hseSBncmF0ZWShYU2iaGNhdGVnb3J5oWFCR2NoZWVzZXNlaXRlbXOhYUyDoWFNo2RpdGVtoWFCU3Bhcm1pZ2lhbm8tcmVnZ2lhbm9mYW1vdW50oWFCTTIgY3VwcyBncmF0ZWRjYWdloWFCSTI0IG1vbnRoc6FhTaNkaXRlbaFhQkdyaWNvdHRhZmFtb3VudKFhQkUyIGxic2R0eXBloWFCSndob2xlIG1pbGuhYU2jZGl0ZW2hYUJKbW96emFyZWxsYWZhbW91bnShYUJNMSBsYiBzaHJlZGRlZGR0eXBloWFCTGxvdy1tb2lzdHVyZWxpbnN0cnVjdGlvbnOhYUyDoWFNpWRzdGVwoWFOYTFldGl0bGWhYUJSUHJlcGFyZSBNZWF0IFNhdWNla2Rlc2NyaXB0aW9uoWFCWD9Ccm93biBwYW5jZXR0YSwgYWRkIGdyb3VuZCBtZWF0cywgY29vayB3aXRoIHZlZ2V0YWJsZXMgYW5kIHdpbmVkdGltZaFhTmI0NWt0ZW1wZXJhdHVyZaFhQkttZWRpdW0taGlnaKFhTaVkc3RlcKFhTmEyZXRpdGxloWFCTU1ha2UgQmVjaGFtZWxrZGVzY3JpcHRpb26hYUJYOkNyZWF0ZSByb3V4IHdpdGggYnV0dGVyIGFuZCBmbG91ciwgZ3JhZHVhbGx5IGFkZCB3YXJtIG1pbGtkdGltZaFhTmIyMGR0aXBzoWFMgqFhQlghV2hpc2sgY29uc3RhbnRseSB0byBwcmV2ZW50IGx1bXBzoWFCWCdLZWVwIG1pbGsgd2FybSBmb3Igc21vb3RoIGluY29ycG9yYXRpb26hYU2lZHN0ZXChYU5hM2V0aXRsZaFhQk5MYXllciBBc3NlbWJseWtkZXNjcmlwdGlvbqFhQlg8QWx0ZXJuYXRlIGxheWVycyBvZiBwYXN0YSwgbWVhdCBzYXVjZSwgYmVjaGFtZWwsIGFuZCBjaGVlc2VzZHRpbWWhYU5iMzBmbGF5ZXJzoWFMg6FhTaJlb3JkZXKhYU5hMWpjb21wb25lbnRzoWJCU4RKbWVhdF9zYXVjZUVwYXN0YUhiZWNoYW1lbEdyaWNvdHRhoWFNomVvcmRlcqFhTmEyamNvbXBvbmVudHOhYkJThEVwYXN0YUptZWF0X3NhdWNlSGJlY2hhbWVsSm1venphcmVsbGGhYU2iZW9yZGVyoWFOYTNqY29tcG9uZW50c6FiQlOERXBhc3RhSm1lYXRfc2F1Y2VIYmVjaGFtZWxKcGFybWlnaWFub2ludXRyaXRpb26hYU2ia3Blcl9zZXJ2aW5noWFNp2hjYWxvcmllc6FhTmM2ODBncHJvdGVpbqFhTmI0Mm1jYXJib2h5ZHJhdGVzoWFOYjM1Y2ZhdKFhTmIzOGVmaWJlcqFhTmEzZnNvZGl1baFhTmQxMjUwa2Nob2xlc3Rlcm9soWFOYzE0NWxkYWlseV92YWx1ZXOhYU2kZ3Byb3RlaW6hYU5iODRpdml0YW1pbl9hoWFOYjI1Z2NhbGNpdW2hYU5iNDVkaXJvbqFhTmIyMGllcXVpcG1lbnShYUyEoWFNomRpdGVtoWFCUDl4MTMgYmFraW5nIGRpc2hpZXNzZW50aWFsoWRCT09M9aFhTaJkaXRlbaFhQk1sYXJnZSBza2lsbGV0aWVzc2VudGlhbKFkQk9PTPWhYU2iZGl0ZW2hYUJOaGVhdnkgc2F1Y2VwYW5pZXNzZW50aWFsoWRCT09M9aFhTaNkaXRlbaFhQk1wYXN0YSBtYWNoaW5laWVzc2VudGlhbKFkQk9PTPRrYWx0ZXJuYXRpdmWhYUJTc3RvcmUtYm91Z2h0IHNoZWV0c2x3aW5lX3BhaXJpbmehYU2jZ3ByaW1hcnmhYUJQQ2hpYW50aSBDbGFzc2ljb2xhbHRlcm5hdGl2ZXOhYkJTg0pTYW5naW92ZXNlTkJhcmJlcmEgZCdBbGJhTU1vbnRlcHVsY2lhbm9sc2VydmluZ190ZW1woWFCSDYwLTY1wrBGZ3N0b3JhZ2WhYU2ibHJlZnJpZ2VyYXRvcqFhTaJoZHVyYXRpb26hYUJIMy00IGRheXNpY29udGFpbmVyoWFCT2NvdmVyZWQgdGlnaHRseWdmcmVlemVyoWFNomhkdXJhdGlvbqFhQkgzIG1vbnRoc2xpbnN0cnVjdGlvbnOhYUyDoWFCWB9Db29sIGNvbXBsZXRlbHkgYmVmb3JlIGZyZWV6aW5noWFCWBlXcmFwIGluIHBsYXN0aWMgdGhlbiBmb2lsoWFCWB5UaGF3IG92ZXJuaWdodCBpbiByZWZyaWdlcmF0b3JncmV2aWV3c6FhTIOhYU2mZnJhdGluZ6FhTmE1Z2NvbW1lbnShYUJYPkFic29sdXRlbHkgaW5jcmVkaWJsZSEgV29ydGggZXZlcnkgbWludXRlIG9mIHByZXBhcmF0aW9uIHRpbWUuaHJldmlld2VyoWFCT2NoZWZfbWFyaW9fMjAyMWRkYXRloWFCSjIwMjEtMTItMTVodmVyaWZpZWShZEJPT0z1bWhlbHBmdWxfdm90ZXOhYU5iNDehYU2mZnJhdGluZ6FhTmE1Z2NvbW1lbnShYUJYW0ZhbWlseSByZWNpcGUgcGVyZmVjdGlvbi4gTWFkZSB0aGlzIGZvciBDaHJpc3RtYXMgZGlubmVyIGFuZCBldmVyeW9uZSBhc2tlZCBmb3IgdGhlIHJlY2lwZSFocmV2aWV3ZXKhYUJKbm9ubmFfcm9zYWRkYXRloWFCSjIwMjEtMTItMjVodmVyaWZpZWShZEJPT0z1bWhlbHBmdWxfdm90ZXOhYU5iMzKhYU2mZnJhdGluZ6FhTmE0Z2NvbW1lbnShYUJYLkRlbGljaW91cyBidXQgdmVyeSB0aW1lIGNvbnN1bWluZy4gUGxhbiBhaGVhZCFocmV2aWV3ZXKhYUJPYnVzeV9wYXJlbnRfMTIzZGRhdGWhYUJKMjAyMS0xMS0yOGh2ZXJpZmllZKFkQk9PTPVtaGVscGZ1bF92b3Rlc6FhTmIxOG5yZWNpcGVfaGlzdG9yeaFhTaRmb3JpZ2luoWFCVUVtaWxpYS1Sb21hZ25hLCBJdGFseXJmYW1pbHlfZ2VuZXJhdGlvbnOhYU5hNG5maXJzdF9yZWNvcmRlZKFhQkQxOTIzbW1vZGlmaWNhdGlvbnOhYUyCoWFNo2R5ZWFyoWFCRDE5NjVmY2hhbmdloWFCWBhBZGRlZCB3aW5lIHRvIG1lYXQgc2F1Y2VmcmVhc29uoWFCVUVuaGFuY2VkIGZsYXZvciBkZXB0aKFhTaNkeWVhcqFhQkQxOTg3ZmNoYW5nZaFhQlgeSW5jcmVhc2VkIGNoZWVzZSBibGVuZCB2YXJpZXR5ZnJlYXNvbqFhQlgaSW1wcm92ZWQgdGV4dHVyZSBhbmQgdGFzdGVsY29va2luZ190aXBzoWFMg6FhTaNoY2F0ZWdvcnmhYUJLcHJlcGFyYXRpb25jdGlwoWFCWDNNYWtlIHNhdWNlIGRheSBiZWZvcmUgZm9yIGJldHRlciBmbGF2b3IgZGV2ZWxvcG1lbnRqaW1wb3J0YW5jZaFhQkRoaWdooWFNo2hjYXRlZ29yeaFhQkhhc3NlbWJseWN0aXChYUJYM0xldCBlYWNoIGxheWVyIGNvb2wgc2xpZ2h0bHkgYmVmb3JlIGFkZGluZyB0aGUgbmV4dGppbXBvcnRhbmNloWFCRm1lZGl1baFhTaNoY2F0ZWdvcnmhYUJGYmFraW5nY3RpcKFhQlg1Q292ZXIgd2l0aCBmb2lsIGZvciBmaXJzdCBob3VyLCB0aGVuIHVuY292ZXIgdG8gYnJvd25qaW1wb3J0YW5jZaFhQkRoaWdocENvbnN1bWVkQ2FwYWNpdHmkaVRhYmxlTmFtZW1wYXN0YS1yZWNpcGVzbUNhcGFjaXR5VW5pdHP7QCEAAAAAAABxUmVhZENhcGFjaXR5VW5pdHP7QCEAAAAAAAByV3JpdGVDYXBhY2l0eVVuaXRzAA==
"""
        tags: ["serde-benchmark"]
    }
])
operation GetItem {
    input: GetItemInput
    output: GetItemOutput
}

structure GetItemInput {
    @required
    TableName: String
    @required
    Key: AttributeValueMap

    AttributesToGet: AttributeNameList
    ConsistentRead: Boolean
    ReturnConsumedCapacity: String
    ProjectionExpression: String
    ExpressionAttributeNames: ExpressionAttributeNameMap
}

structure GetItemOutput {
    Item: AttributeValueMap
    ConsumedCapacity: ConsumedCapacity
}

list AttributeNameList {
    member: String
}

map ExpressionAttributeNameMap {
    key: String
    value: String
}

structure ConsumedCapacity {
    TableName: String
    CapacityUnits: Double
    ReadCapacityUnits: Double
    WriteCapacityUnits: Double
}
