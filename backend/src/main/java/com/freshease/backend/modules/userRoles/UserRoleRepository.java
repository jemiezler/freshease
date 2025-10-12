/*
 * Click nbfs://nbhost/SystemFileSystem/Templates/Licenses/license-default.txt to change this license
 * Click nbfs://nbhost/SystemFileSystem/Templates/Classes/Interface.java to edit this template
 */

package com.freshease.backend.modules.userRoles;

import java.util.UUID;

import org.springframework.data.jpa.repository.JpaRepository;

/**
 *
 * @author jemiezler
 */
public interface UserRoleRepository extends JpaRepository<UserRoleEntity, UUID> {
    boolean existsByUser_IdAndRole_Id(UUID userId, UUID roleId);
}
