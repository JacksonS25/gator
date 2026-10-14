# Gator
This is your new RSS Feed Aggregator in the Command Line, follow the feeds you want, store them to be read later, and browse them.

## Motivation
It's nice to stay up to date on the latest news fo ryour favorite websites. This is where the gator comes in!
It allows you to get your favorite websites RSS Feeds stored in a database to be viewed at will.
This also allows for easy access to the most recent posts.

## Quick Start
1. You can clone the repo and make sure that Go and Postgres are installed. If not, use:
   ```wget https://go.dev
   sudo apt update
   sudo apt install postgresql postgresql-contrib```
3. To use the gator CLI use ```go install``` from the root directory of the repository.
4. To make sure the database is setup:
   Set your password for postgres, you can use postgres for simplicity:
   ```sudo passwd postgres```
   Then start the service:
   ```sudo service postgresql start```
   Create a database, you can name it gator if you like:
   ```CREATE DATABASE gator;```
   Connect to the database:
   ```\c gator```
   Alter the databases password, below postgres is an example:
   ```ALTER USER postgres PASSWORD 'postgres';```
   From here you can run SQL queries against the database. Type ```exit``` to leave the psql shell.
5. Next in your home directory create a config.json file name it ".gatorconfig.json". Include this:
   ```{"db_url":"postgres://postgres:postgres@localhost:5432/gator?sslmode=disable","current_user_name":""}```
   Change the names of the user and passwords as necessary.

## Commands
```gator <command_name> <arguments>```
- register <username>: Add user to the users table and login as the user.
- login <username>: Login as that user.
- reset <reset>: Reset the database.
- users: List users.
- agg <duration e.g. 1s or 10m>: Scrape feeds you are following on an interval and adding their posts to the database.
- addfeed <name> <url>: Create a new feed and set the current user to following it.
- feeds: List all registered feeds.
- follow <url>: Follow the feed matching the url.
- following: List the feeds the current user is following.
- unfollow <url>: Unfollow the feed matching the url.
- browse <optional: limit; default: 2>: Load the most recent posts up to the limit.

## Future Features
- User passwords
- Ability to follow and unfollow using the name of the feed
- Additional Housekeeping
- More In-Depth and Varied Browsing

## 🤝 Contributing

### Clone the repo

```bash
git clone https://github.com/JacksonS25/gator
cd gator
```

### Submit a pull request

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.
