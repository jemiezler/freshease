/*
 * Click nbfs://nbhost/SystemFileSystem/Templates/Licenses/license-default.txt to change this license
 * Click nbfs://nbhost/SystemFileSystem/Templates/Classes/Class.java to edit this template
 */

package com.freshease.backend.modules.rolePermissions;

import java.util.UUID;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;

/**
 *
 * @author jemiezler
 */
public interface RolePermissionRepository extends JpaRepository<RolePermissionEntity, UUID> {
  boolean existsByRole_IdAndPermission_Id(UUID roleId, UUID permId);
  @Query("""
    select count(rp) > 0
    from RolePermissionEntity rp
    join UserRoleEntity ur on ur.role = rp.role
    where ur.user.id = :userId and rp.permission.id = :permissionId
  """)
  boolean userHasPermission(UUID userId, UUID permissionId);
}