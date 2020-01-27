task :default => [:help]

def settings_file_path
    current_dir = File.expand_path(File.dirname(__FILE__))
    settings_path = File.join(current_dir, './.settings')
    File.expand_path(settings_path)
end

def encrypted_file
    unless File.exists? settings_file_path
        raise "Settings file not found at '#{settings_file_path}'. Make sure to run rake setup"
    end

    File.read(settings_file_path).strip
end

def encrypted_file_configured?
    return false unless File.exists?(settings_file_path)

    return !File.read(settings_file_path).strip.empty?
end

def encrypted_file_bck
    "#{encrypted_file}.bck"
end

def encrypted_file_parent
    File.dirname encrypted_file
end

def decrypted_folder_parent
    File.dirname decrypted_folder_full_path
end

def decrypted_folder_full_path
    File.expand_path "~/decrypted_folder/#{decrypted_folder_name}"
end

def decrypted_folder_name
    "t/"
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
    puts "settings_file_path           => #{settings_file_path}"

    if !encrypted_file_configured?
        puts "Settings not configured yet. Run rake setup"
    else
        puts "encrypted_file               => #{encrypted_file}"
        puts "encrypted_file_parent        => #{encrypted_file_parent}"
        puts "decrypted_folder_parent      => #{decrypted_folder_parent}"
        puts "decrypted_folder_full_path   => #{decrypted_folder_full_path}"
    end
end

task :e => :encrypt
task :encrypt => [:create_decrypted_folder_parent, :create_encrypted_file_folder, :create_encrypted_file_backup] do
    Dir.chdir(decrypted_folder_parent){        
        sh "zip -er #{encrypted_file} #{decrypted_folder_name}"
    }

    Rake::Task[:clean].invoke
end

task :create_encrypted_file_backup do
    if File.file? encrypted_file
        mv encrypted_file, encrypted_file_bck
    end
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

    if encrypted_file_configured?
        if File.file? encrypted_file_bck
            rm encrypted_file_bck
        end
    end
end

task :setup => :configure_path do
    ["clean", "decrypt", "encrypt"].each { |task_name|
        generate_task_alias_wrapper task_name
    }
end

task :configure_path do
    puts 'Enter the path of your encryped zip file'
    file_path_raw = STDIN.gets.strip
    file_path = File.expand_path file_path_raw
    unless File.exists? file_path
        fail "File does not exists at '#{file_path}'"
    end

    write_file(settings_file_path, file_path)
end

def generate_task_alias_wrapper(task_name)
    generate_task_alias "#{task_name}.sh.command", get_file_contents(task_name)
end

def get_file_contents(task_name)
    %{
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
    File.expand_path "."
end

def get_full_path(filename)
    File.join(basedir, filename)
end

def write_file(file_path, file_content)
    File.open(file_path, "w") do |f|
        f.puts file_content
    end
end
