/*
 * Click nbfs://nbhost/SystemFileSystem/Templates/Licenses/license-default.txt to change this license
 * Click nbfs://nbhost/SystemFileSystem/Templates/Classes/Class.java to edit this template
 */

package com.freshease.backend.modules.roles;

import java.util.UUID;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import com.freshease.backend.modules.rolePermissions.RolePermissionService;

import lombok.RequiredArgsConstructor;

/**
 *
 * @author jemiezler
 */
@RestController
@RequestMapping("/roles")
@RequiredArgsConstructor
public class RoleController {
    private final RoleService roleService;
    private final RolePermissionService rolePermissionService;

    @PostMapping()
    public RoleEntity createRole(@RequestParam String name, @RequestParam(required = false) String description) {
        return roleService.createRole(
                new RoleEntity(
                        java.util.UUID.randomUUID(),
                        name,
                        description,
                        java.time.OffsetDateTime.now(),
                        java.time.OffsetDateTime.now()));
    }

    @GetMapping()
    public String findAll() {
        return roleService.findAll().toString();
    }

    @PostMapping("/roles/{roleId}/grant/{permId}")
    public void grant(@PathVariable UUID roleId, @PathVariable UUID permId) {
        rolePermissionService.grantPermissionToRole(roleId, permId);
    }

}
