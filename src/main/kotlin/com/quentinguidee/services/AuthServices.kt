package com.quentinguidee.services

import com.quentinguidee.client
import com.quentinguidee.dao.gitHubUsersDAO
import com.quentinguidee.dao.sessionsDAO
import com.quentinguidee.dao.usersDAO
import com.quentinguidee.models.db.Session
import io.ktor.client.call.*
import io.ktor.client.request.*
import io.ktor.http.*
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import org.jetbrains.exposed.sql.transactions.transaction

@Serializable
data class GitHubUserBody(
    val email: String,
    val name: String,
    @SerialName("avatar_url")
    val avatarURL: String,
    val login: String,
)

class AuthServices {
    suspend fun githubUser(username: String) = transaction {
        gitHubUsersDAO.get(username)
    }

    suspend fun fetchGitHubUser(token: String): GitHubUserBody {
        return client
            .request("https://api.github.com/user") {
                headers {
                    append(HttpHeaders.Authorization, "token $token")
                }
            }
            .body()
    }

    private suspend fun createAccount(username: String, name: String, email: String, profilePicture: String): Session =
        transaction {
            val user = usersDAO.create(
                username = username,
                name = name,
                email = email,
                profilePicture = profilePicture,
            )

            gitHubUsersDAO.create(user.id, username)

            return@transaction sessionsDAO.create(user.id)
        }

    suspend fun createAccount(gitHubUser: GitHubUserBody) = authServices.createAccount(
        username = gitHubUser.login,
        name = gitHubUser.name,
        email = gitHubUser.email,
        profilePicture = gitHubUser.avatarURL,
    )

    suspend fun session(username: String) = transaction {
        val user = usersDAO.get(username)
        return@transaction sessionsDAO.get(user.id)
    }
}

val authServices = AuthServices()