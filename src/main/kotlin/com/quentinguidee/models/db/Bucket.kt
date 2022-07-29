package com.quentinguidee.models.db

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import org.jetbrains.exposed.dao.id.UUIDTable

enum class BucketType {
    USER_BUCKET
}

object Buckets : UUIDTable() {
    val name = varchar("name", 255)
    val type = enumerationByName("type", 63, BucketType::class)
    val size = integer("size").default(0)
    val maxSize = integer("max_size").nullable()
}

@Serializable
data class Bucket(
    val uuid: String,
    val name: String,
    val type: BucketType,
    val size: Int,
    @SerialName("max_size")
    val maxSize: Int?,
)
