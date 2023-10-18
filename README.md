# SpotiBio


**SpotiBio** shows your current listening music from Spotify in your Twitter bio.

 The name `SpotiBio` comes from Spotify + Bio Twitter.
 
![spotiBio](https://github.com/pooriaanv/spotibio/assets/39620620/6fd4d432-7e58-4824-b38c-720c10b47d5c)

## Requirements
- Docker and Docker Compose
- Twitter and Spotify API credentials.

## Getting started

To run and build SpotiBio, you have to do the following steps:

- Pull the codes
- Move to the project root folder.
- Create a .env file from .env.example.
- Insert your Twitter and Spotify API credentials into the .env file
- Run the following command in the terminal.

```bash
bash start.sh
```
There are two containers, `spotibio-app` and `spotibio-cron`.

`spotibio-cron` handles cronjob and executes the script **every 2 minutes**.
## More Options

#### Default Twitter bio:
If you want to keep your own Twitter bio beside the names of the songs.

 Set `DEFAUL_DESCRIPTION` value in `.env` according to your bio.

#### Change cronjob schedule:
To change the script's schedule, change `cronfile` in the cron folder. 


> **Note** After the first run and build, for every change that you apply, you should rebuild by following the command:
```bash
bash start.sh rebuild
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)
