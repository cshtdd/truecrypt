task :default => [:help]

def encrypted_file
    return File.expand_path "~/Documents/t/t.zip"
end

def encrypted_file_parent
    return File.dirname encrypted_file
end

def decrypted_folder_parent
    return File.dirname decrypted_folder_full_path
end

def decrypted_folder_full_path
    return File.expand_path "~/decrypted_folder/#{decrypted_folder_name}"
end

def decrypted_folder_name
    return "t/"
end

task :help do
    puts "TrueCrypt replacement"

    puts ""
    puts "Tasks"
    puts "="*16
    puts "decrypt (d)"
    puts "encrypt (e)"
    puts "clean   (c)"
    puts "setup"

    puts ""
    puts "Settings"
    puts "="*16
    puts "encrypted_file               => #{encrypted_file}"
    puts "encrypted_file_parent        => #{encrypted_file_parent}"
    puts "decrypted_folder_parent      => #{decrypted_folder_parent}"
    puts "decrypted_folder_full_path   => #{decrypted_folder_full_path}"
end

task :e => :encrypt
task :encrypt => [:create_decrypted_folder_parent, :create_encrypted_file_folder] do
    encrypted_file_bck = "#{encrypted_file}.bck"
    if File.file? encrypted_file
        mv encrypted_file, encrypted_file_bck
    end

    Dir.chdir(decrypted_folder_parent){        
        sh "zip -er #{encrypted_file} #{decrypted_folder_name}"
    }

    if File.file? encrypted_file_bck
        rm encrypted_file_bck
    end

    Rake::Task[:clean].invoke
end

task :create_decrypted_folder_parent do
    `mkdir #{decrypted_folder_parent}`
end

task :create_encrypted_file_folder do
    `mkdir #{encrypted_file_parent}`
end

task :d => :decrypt
task :decrypt => [:ensure_decrypted_folder_doesnt_exist, :create_decrypted_folder_parent] do
    sh "unzip -n #{encrypted_file} -d #{decrypted_folder_parent}"
end

task :ensure_decrypted_folder_doesnt_exist do
    if File.exists? decrypted_folder_full_path
        fail "Decrypted folder already exists at '#{decrypted_folder_full_path}' please manually delete it before continuing"
    end
end

task :c => :clean
task :clean do
    rm_rf decrypted_folder_full_path
end

task :setup do
    ["clean", "decrypt", "encrypt"].each { |task_name|
        generate_task_alias_wrapper task_name
    }
end

def generate_task_alias_wrapper(task_name)
    generate_task_alias "#{task_name}.sh.command", get_file_contents(task_name)
end

def get_file_contents(task_name)
    return %{
#!/bin/sh
pushd #{basedir}
rake #{task_name}
read -p "Press any key to continue..."
}
end

def generate_task_alias(filename, file_contents)
    file_path = get_full_path filename
    write_file file_path, file_contents
    chmod "+x", file_path
end

def basedir
   return File.expand_path "."
end

def get_full_path(filename)
    return File.join(basedir, filename)
end

def write_file(file_path, file_content)
    File.open(file_path, "w") do |f|
        f.puts file_content
    end
end
