 
 - add dynamic linking within the idea files like logseq has, this should:
    - [[1892]] - this will link to idea number 1892 
    - when the cursor is over 1892, need to build vim scripts so that `gd` will 
      open up that idea in a new tab
    - potentially add a command which links these together in the filenames?
      something like:
      ideas/a,000419,2020-02-08,e2020-02-08,aib,communications,admin,REF:1892
       - could also just perform this task for _all_ idea files with a
         big-scan/rename operation, probably easier
 - "neuron" view browser mode:
                                                                           
   ┌─1892────────────────┐                                                             
   │tag1, footag,  artag │                                    
   ├─────────────────────┤                                                                       
   │... exurpt from this │
   │document             │
   │                 ... │
   └─────────────────────┘ 
                                                                                        
   ┌─1892────────────────┐        ┌─1892────────────────┐        
   │tag1, footag,  artag │        │tag1, footag,  artag │
   ├─────────────────────┤        ├─────────────────────┤
   │... exurpt from this │        │.CENTER            s │    
   │document             │        │       t             │    
   │                 ... │        │                 ... │    
   └─────────────────────┘        └─────────────────────┘
                                                                                        
   ┌─1892────────────────┐        ┌─1892────────────────┐                
   │tag1, footag,  artag │        │tag1, footag,  artag │ 
   ├─────────────────────┤        ├─────────────────────┤
   │... exurpt from this │        │... exurpt from this │            
   │document             │        │document             │            
   │                 ... │        │                 ... │            
   └─────────────────────┘        └─────────────────────┘
                                                                                        
   ┌─1892────────────────┐  ┌─1892──────────┐                      
   │tag1, footag,  artag │  │tag1, footag,  │       
   ├─────────────────────┤  ├───────────────┤      
   │... exurpt from this │  │... exurpt fr  │                  
   │document             │  │document       │                  
   │                 ... │  │               │                  
   └─────────────────────┘  └───────────────┘      
      - have to generate a score matrix of relatability of the various idea
        - common tags
        - explicit connections linking ideas
        - overlapping usage of words? 
           - larger score ...5 words in a row
           - next         ...4 words in a row
           - next         ...3 words in a row
           - next         ...2 words in a row
           - next         ...single word (lowest score) 
         - OR just look at the greatest size of overlapping characters, 
           don't even consider words, just consider chunks of overlapping
           character sets
      - problem. Don't want the map to be disorienting too much 
        when browsing, but naturally the elements in the map
        will need to shift which each step within the map. 
          SOLUTION: a brief animation of the elements as they shift 
          around from each other as the ideas are being input. 

      THREE CONCENTRIC RINGS
      - center is the largest
      - next can be the next largest
      - 3rd can be the smallest

      VIEW 1
       - just a single layer of connections around a single node
         - using strong detail
      VIEW 2 
       - expanded single layer of connections, weaker connections shown
         but everything has less detail
      VIEW 3
       - two layers of connection, more-elements = less detail 

EACH connection to have the connectivity score associated with it. 

ZOOM functionality - ability to select an element from the CUI and enlarge it
therefore increasing its viewable content space. DO NOT reorient other boxes, 
just overlap them from the box being enlarged. 
                                                                          
                                                                          
                                                                          
                                                                          
 - LOW PRIORITY                                                                          
   - view trash
   - recover from trash
                                                                          
                                                                          
                                                                          
                                                                          
                                                                          
                                                                          
                                                                          
                                                                          
