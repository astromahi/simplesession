simplesession
=============

File based session support for Golang

### Features

-   Simple, precise & clean code
-   Fast & efficient
-   No dependency injection
-   Using go standard library only (No third-party library included)
-   Best suit for the applications need file based session handling

### Usage

1.  Create New Session
    
    New creates a fresh session environment

        option := &simplesession.Option{
            Path:     "/",
            Domain:   "example.com", //your domain
            MaxAge:   24 * 60 * 60, //expiry time of session cookie
            Secure:   false,
            HttpOnly: true,
        }

        session, err := simplesession.New(res, option)
        if err != nil {
            return err
        }

        // Do something with session

2.  Set Session Data
    
    simplesession uses map[string]interface{} as its storage for handling data at runtime which
    greatly improves the efficiency

        session.Set(key, val)

        Ex.
        session.Set("id", 101)
        
3.  Write Session

    Write writes the session data to file for persistent and later use.  
    It takes http.ResponseWriter, http.Request as parameter

        err := session.Write(res, req)
        if err != nil {
            return err
        }
        return nil

4.  Read Session
    
    Read reads the data from stored session.  
    It takes http.Request and gives stored session
    
        session, err := simplesession.Read(req)
        if err != nil {
            return err
        }
        
        // Do something with session

5.  Get Session Data
    
    Get returns the stored data from read session

        data := session.Get(key)

        Ex.
        userId := session.Get("id")

6.  Delete Session Data
    
    Del deletes the user session data from session

        session.Del("id")

    After deletion, you should save the modified session

        if err := session.Write(res, req); err != nil {
            return err
        }
        return nil

7.  Destroy Session

    Destroy destroys the whole session & removes the session file from disk.

        err := session.Destroy(res)
        if err != nil {
            return err
        }
        return nil
 
