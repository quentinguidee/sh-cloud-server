package com.quentinguidee.models.db

import kotlinx.serialization.json.JsonObject
import kotlinx.serialization.json.buildJsonObject
import kotlinx.serialization.json.put
import org.jetbrains.exposed.dao.UUIDEntity
import org.jetbrains.exposed.dao.UUIDEntityClass
import org.jetbrains.exposed.dao.id.EntityID
import org.jetbrains.exposed.dao.id.UUIDTable
import java.util.*

enum class BucketType {
    USER_BUCKET
}

object Buckets : UUIDTable() {
    val name = varchar("name", 255)
    val type = enumerationByName("type", 63, BucketType::class)
    val size = integer("size").default(0)
    val maxSize = integer("max_size").nullable()
}

class Bucket(id: EntityID<UUID>) : UUIDEntity(id) {
    companion object : UUIDEntityClass<Bucket>(Buckets)

    var name by Buckets.name
    var type by Buckets.type
    var size by Buckets.size
    var maxSize by Buckets.maxSize

    var users by User via UserBuckets
    var rootNode: Node? = null

    fun toJSON(): JsonObject {
        return buildJsonObject {
            put("name", name)
            put("type", type.name)
            put("size", size)
            put("max_size", maxSize)
            if (rootNode != null)
                put("root_node", rootNode!!.toJSON())
        }
    }
}