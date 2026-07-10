# image_service
    Common external project entities. Service actions 

## - Actions: 

### -- 1) Add image  
**input:**
- oauth_token
- name (string)
- image (form data)
**output:**
- image_id
- is success


## -- 2) Get image
**input:**
- oauth_token
- id
**output:**
- image media

## -- 3) block image 
**input:**
- oauth_token
- image_id 
**output:**
- is success

## -- 4) approved image
**input:**
- oauth_token
- image_id 
**output:**
- is success



## - Models: 
### -- 1) image  
- id (int)
- media_type (string)
- name (string)
- created_at (float64)
- updated_at (float64)
- status (string)(blocked/approved/unmoderated)

## -- 2) client
- id
- created_at (float64)
- name (string)
- token (string)
- expired_at (float64)
- login (string)
- pass_hash (string)
