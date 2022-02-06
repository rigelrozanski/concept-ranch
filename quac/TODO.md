 
 - fix various bugs preventing error output
 - search or list by file contents as well as tags
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

                consumes: ─────┐                     common subtags       
                               │                             │            
                               │     ┌────────────first tag──┘            
                               │     │                                    
                               │     │                                    
                      ┌─1892───┴─────┴──────────┐                             
                      │  some idea blah blah    ├────────tag-2
                      │  blah blah blah blah    │
                      │                         │                         
                      │                         │                         
                      └─────────────────────────┘                         
                                                                          
                                                                          
                                                                          
                                                                          
                                                                          
 - LOW PRIORITY                                                                          
   - view trash
   - recover from trash
                                                                          
                                                                          
                                                                          
                                                                          
                                                                          
                                                                          
                                                                          
                                                                          
